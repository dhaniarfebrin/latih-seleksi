// Package persentase generates chained-percentage TIU questions
// (pattern TIU.numerik.persentase-berantai).
package persentase

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.numerik.persentase-berantai"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
)

// Misconception labels attached to wrong options.
const (
	MisJumlahLangsung   = "jumlah-langsung"
	MisLangkahTerakhir  = "langkah-terakhir-saja"
	MisTandaTerbalik    = "tanda-terbalik"
	MisSalahHitungDekat = "salah-hitung-dekat"
)

// params are the numbers behind one question: changes holds the signed percent
// of each step (positive = increase) and context selects the sentence template.
type params struct {
	changes []int
	context int
}

// answer is the exact result of solve: the net change in percent is num/den.
type answer struct {
	num  int64
	den  int64
	text string
}

type contextTemplate struct {
	subject string // opening noun phrase
	final   string // label of the end value
	ref     string // what the last change is measured against
}

var contexts = []contextTemplate{
	{"Harga suatu barang", "harga akhir", "harga"},
	{"Nilai investasi seorang pegawai", "nilai akhir", "nilai"},
	{"Harga sebuah tiket pesawat", "harga akhir", "harga"},
	{"Gaji pokok seorang karyawan", "gaji akhir", "gaji"},
	{"Tarif parkir sebuah pusat perbelanjaan", "tarif akhir", "tarif"},
	{"Harga sebuah buku pelajaran", "harga akhir", "harga"},
	{"Upah harian seorang pekerja", "upah akhir", "upah"},
	{"Iuran bulanan sebuah perumahan", "iuran akhir", "iuran"},
	{"Biaya pendaftaran sebuah kursus", "biaya akhir", "biaya"},
	{"Target penjualan sebuah toko", "target akhir", "target"},
}

// Family implements gen.Family for chained percentages.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Persentase berantai" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisJumlahLangsung:   "Menjumlahkan persentase langsung, padahal tiap perubahan berlaku pada nilai terbaru",
		MisLangkahTerakhir:  "Hanya menghitung perubahan terakhir terhadap nilai awal",
		MisTandaTerbalik:    "Tanda kenaikan dan penurunan tertukar",
		MisSalahHitungDekat: "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Parameters whose net
// change has more than one decimal are rejected; attempt k uses a seed
// derived from (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, oneDecimal)
}

func oneDecimal(p params) bool {
	a, _ := solve(p)
	return a.num*10%a.den == 0
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("persentase: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("persentase: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

func pickParams(rng *rand.Rand, band gen.Band) params {
	n := 2
	if band == gen.Band3 {
		n = 3
	}
	changes := make([]int, n)
	for i := range changes {
		switch band {
		case gen.Band1:
			changes[i] = (rng.Intn(12) + 1) * 5
		case gen.Band2:
			changes[i] = rng.Intn(50) + 1
		default:
			// multiples of 5 keep the 3-step key at one decimal often enough
			changes[i] = (rng.Intn(10) + 1) * 5
		}
	}
	signs := make([]bool, n) // true = increase
	for i := range signs {
		signs[i] = rng.Intn(2) == 0
	}
	if n == 2 {
		signs[1] = !signs[0]
	} else if signs[0] == signs[1] && signs[1] == signs[2] {
		signs[2] = !signs[2]
	}
	for i, up := range signs {
		if !up {
			changes[i] = -changes[i]
		}
	}
	return params{changes: changes, context: rng.Intn(len(contexts))}
}

func solve(p params) (answer, []string) {
	steps := []string{"Misalkan nilai awal = 100."}
	num, digits := int64(100), 0 // current value = num / 10^digits
	for i, c := range p.changes {
		num *= int64(100 + c)
		digits += 2
		steps = append(steps, fmt.Sprintf("Langkah %d: %s %d%%, nilai x %s = %s.",
			i+1, verb(c), abs(c), fmtFrac(int64(100+c), 2, false), fmtFrac(num, digits, true)))
	}
	den := pow10(digits)
	net := num - 100*den
	ans := answer{num: net, den: den, text: percentText(net, digits)}
	steps = append(steps, fmt.Sprintf("Perubahan bersih = %s - 100 = %s, yaitu %s.",
		fmtFrac(num, digits, true), signedFrac(net, digits), ans.text))
	return ans, steps
}

func build(p params, seed int64) gen.Draft {
	ans, steps := solve(p)
	keyTenths := int(ans.num * 10 / ans.den)

	type cand struct {
		tenths int
		mis    string
	}
	sum, last := 0, p.changes[len(p.changes)-1]
	for _, c := range p.changes {
		sum += c
	}
	cands := []cand{
		{sum * 10, MisJumlahLangsung},
		{last * 10, MisLangkahTerakhir},
		{-keyTenths, MisTandaTerbalik},
	}
	for d := 5; d <= 100; d += 5 {
		cands = append(cands, cand{keyTenths + d, MisSalahHitungDekat}, cand{keyTenths - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	opts := []opt{{ans.text, ""}}
	seen := map[string]bool{ans.text: true}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		t := percentText(int64(c.tenths), 1)
		if seen[t] {
			continue
		}
		seen[t] = true
		opts = append(opts, opt{t, c.mis})
	}

	rng := rand.New(rand.NewSource(derive(seed, shuffleSalt)))
	rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	options := make([]gen.Option, len(opts))
	answerKey := ""
	for i, o := range opts {
		key := string(rune('A' + i))
		options[i] = gen.Option{Key: key, Text: o.text, MisconceptionID: o.mis}
		if o.mis == "" {
			answerKey = key
		}
	}
	return gen.Draft{
		PatternID: PatternID,
		Stem:      stem(p),
		Options:   options,
		AnswerKey: answerKey,
		Steps:     steps,
		Seed:      seed,
	}
}

func stem(p params) string {
	t := contexts[p.context%len(contexts)]
	var b strings.Builder
	b.WriteString(t.subject)
	for i, c := range p.changes {
		switch i {
		case 0:
			b.WriteString(" ")
		case 1:
			b.WriteString(", kemudian ")
		default:
			b.WriteString(", lalu ")
		}
		if c > 0 {
			fmt.Fprintf(&b, "dinaikkan %d%%", c)
		} else {
			fmt.Fprintf(&b, "diturunkan %d%%", -c)
		}
		if i > 0 {
			prev := "kenaikan"
			if p.changes[i-1] < 0 {
				prev = "penurunan"
			}
			fmt.Fprintf(&b, " dari %s setelah %s", t.ref, prev)
		}
	}
	fmt.Fprintf(&b, ". Dibandingkan %s semula, %s ...", t.ref, t.final)
	return b.String()
}

func verb(c int) string {
	if c < 0 {
		return "turun"
	}
	return "naik"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func pow10(n int) int64 {
	r := int64(1)
	for ; n > 0; n-- {
		r *= 10
	}
	return r
}

// fmtFrac formats num / 10^digits with a decimal comma.
func fmtFrac(num int64, digits int, trim bool) string {
	s := strconv.FormatInt(num, 10)
	for len(s) <= digits {
		s = "0" + s
	}
	whole, frac := s[:len(s)-digits], s[len(s)-digits:]
	if trim {
		frac = strings.TrimRight(frac, "0")
	}
	if frac == "" {
		return whole
	}
	return whole + "," + frac
}

func signedFrac(num int64, digits int) string {
	if num < 0 {
		return "-" + fmtFrac(-num, digits, true)
	}
	return fmtFrac(num, digits, true)
}

// percentText renders num / 10^digits percent as "naik X%", "turun X%" or "tetap".
func percentText(num int64, digits int) string {
	switch {
	case num == 0:
		return "tetap"
	case num > 0:
		return "naik " + fmtFrac(num, digits, true) + "%"
	default:
		return "turun " + fmtFrac(-num, digits, true) + "%"
	}
}
