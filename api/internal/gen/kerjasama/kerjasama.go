// Package kerjasama generates work-rate TIU questions
// (pattern TIU.numerik.kerja-sama).
package kerjasama

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.numerik.kerja-sama"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
	maxDays     = 120 // upper bound of any answer
	maxTime     = 90  // upper bound of any time given in the stem
)

// Misconception labels attached to wrong options.
const (
	MisRataRataWaktu    = "rata-rata-waktu"
	MisJumlahWaktu      = "jumlah-waktu"
	MisWaktuTercepat    = "waktu-tercepat"
	MisHanyaSisaHari    = "hanya-sisa-hari"
	MisWaktuSendiri     = "waktu-sendiri"
	MisSalahHitungDekat = "salah-hitung-dekat"
)

// params are the numbers behind one question: times holds the days each worker
// needs alone; together is the number of days the first two work jointly
// before the second stops (0 when everyone works until the job is done).
type params struct {
	times    []int
	together int
	context  int
}

type contextTemplate struct {
	job string // gerund phrase naming the job
}

var contexts = []contextTemplate{
	{"mengecat sebuah pagar"},
	{"membersihkan sebuah gudang"},
	{"menyusun sebuah laporan"},
	{"menggali sebuah saluran air"},
	{"memasang keramik sebuah ruangan"},
	{"memanen sebidang ladang jagung"},
	{"menjahit seragam sekolah satu kelas"},
	{"mencetak setumpuk brosur"},
	{"merapikan arsip sebuah kantor"},
	{"membersihkan sebuah kolam"},
	{"memperbaiki sebuah jalan desa"},
	{"mengecat sebuah ruang kelas"},
}

var names = []string{"A", "B", "C"}

// Family implements gen.Family for work-rate questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Kerja sama" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisRataRataWaktu:    "Merata-ratakan waktu masing-masing, padahal bekerja bersama pasti lebih cepat dari yang tercepat",
		MisJumlahWaktu:      "Menjumlahkan waktu masing-masing, padahal yang dijumlahkan adalah laju kerja",
		MisWaktuTercepat:    "Mengambil waktu pekerja tercepat, mengabaikan bantuan pekerja lain",
		MisHanyaSisaHari:    "Hanya menghitung hari kerja sendirian, lupa menambahkan hari kerja bersama",
		MisWaktuSendiri:     "Memakai waktu penyelesaian sendirian, mengabaikan pekerjaan yang sudah dibantu",
		MisSalahHitungDekat: "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Parameters whose total is
// not a whole number of days are rejected; attempt k uses a seed derived from
// (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, wholeDays)
}

func wholeDays(p params) bool {
	for _, t := range p.times {
		if t > maxTime {
			return false
		}
	}
	d, _ := solve(p)
	return d != nil && d.IsInt() && d.Sign() > 0 && d.Num().Int64() <= maxDays
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("kerjasama: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("kerjasama: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

func pickParams(rng *rand.Rand, band gen.Band) params {
	ctx := rng.Intn(len(contexts))
	switch band {
	case gen.Band1:
		// a = k*u*(u+v), b = k*v*(u+v) gives ab/(a+b) = k*u*v, a whole number;
		// k is capped so both times stay within maxTime
		u, v := rng.Intn(9)+1, rng.Intn(9)+1
		hi := u
		if v > hi {
			hi = v
		}
		kMax := maxTime / (hi * (u + v))
		if kMax < 1 {
			return params{context: ctx} // rejected by solve
		}
		k := rng.Intn(kMax) + 1
		return params{times: []int{k * u * (u + v), k * v * (u + v)}, context: ctx}
	case gen.Band2:
		// times n/x, n/y, n/z with n a multiple of lcm(x,y,z,x+y+z) make the
		// joint time n/(x+y+z) whole
		x, y, z := rng.Intn(6)+1, rng.Intn(6)+1, rng.Intn(6)+1
		n := lcm(lcm(x, y), lcm(z, x+y+z)) * (rng.Intn(3) + 1)
		if n > maxDays || n/(x+y+z) < 2 {
			return params{context: ctx} // rejected by solve
		}
		return params{times: []int{n / x, n / y, n / z}, context: ctx}
	default:
		// B stops after k days together; k is a multiple of b/gcd(a,b) so that
		// A's remaining time is whole, and k must leave work for A.
		a, b := rng.Intn(35)+6, rng.Intn(35)+6
		step := b / gcd(a, b)
		limit := (a*b - 1) / (a + b) // largest k with k*(a+b) < a*b
		if a == b || step > limit {
			return params{context: ctx}
		}
		k := step * (rng.Intn(limit/step) + 1)
		return params{times: []int{a, b}, together: k, context: ctx}
	}
}

// solve returns the total days (nil when params are malformed) and the steps.
func solve(p params) (*big.Rat, []string) {
	if len(p.times) < 2 || p.together < 0 {
		return nil, nil
	}
	rate := new(big.Rat)
	parts := make([]string, 0, len(p.times))
	rates := make([]string, 0, len(p.times))
	var steps []string
	for i, t := range p.times {
		if t <= 0 {
			return nil, nil
		}
		r := big.NewRat(1, int64(t))
		rate.Add(rate, r)
		parts = append(parts, r.RatString())
		rates = append(rates, fmt.Sprintf("laju %s = 1/%d", names[i], t))
	}
	steps = append(steps,
		"Anggap seluruh pekerjaan = 1. "+capitalize(strings.Join(rates, ", "))+" pekerjaan per hari.",
		fmt.Sprintf("Laju bersama = %s = %s pekerjaan per hari.", strings.Join(parts, " + "), rate.RatString()))
	if p.together == 0 {
		days := new(big.Rat).Inv(rate)
		steps = append(steps, fmt.Sprintf("Waktu bersama = 1 : %s = %s hari.", rate.RatString(), days.RatString()))
		return days, steps
	}
	k := big.NewRat(int64(p.together), 1)
	done := new(big.Rat).Mul(k, rate)
	rest := new(big.Rat).Sub(big.NewRat(1, 1), done)
	if rest.Sign() <= 0 {
		return nil, nil
	}
	solo := new(big.Rat).Mul(rest, big.NewRat(int64(p.times[0]), 1))
	total := new(big.Rat).Add(k, solo)
	steps = append(steps,
		fmt.Sprintf("Dalam %d hari bersama selesai %d x %s = %s pekerjaan, sisa %s.", p.together, p.together, rate.RatString(), done.RatString(), rest.RatString()),
		fmt.Sprintf("A sendirian menyelesaikan sisa itu dalam %s x %d = %s hari.", rest.RatString(), p.times[0], solo.RatString()),
		fmt.Sprintf("Total = %d + %s = %s hari.", p.together, solo.RatString(), total.RatString()))
	return total, steps
}

func build(p params, seed int64) gen.Draft {
	total, steps := solve(p)
	key := int(total.Num().Int64())

	type cand struct {
		days int
		mis  string
	}
	var cands []cand
	if p.together == 0 {
		sum, min := 0, p.times[0]
		for _, t := range p.times {
			sum += t
			if t < min {
				min = t
			}
		}
		if sum%len(p.times) == 0 {
			cands = append(cands, cand{sum / len(p.times), MisRataRataWaktu})
		}
		cands = append(cands, cand{sum, MisJumlahWaktu}, cand{min, MisWaktuTercepat})
	} else {
		cands = append(cands, cand{key - p.together, MisHanyaSisaHari}, cand{p.times[0], MisWaktuSendiri})
	}
	for d := 1; d <= maxDays; d++ {
		cands = append(cands, cand{key + d, MisSalahHitungDekat}, cand{key - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	opts := []opt{{daysText(key), ""}}
	seen := map[string]bool{daysText(key): true}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		if c.days <= 0 {
			continue
		}
		t := daysText(c.days)
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
		k := string(rune('A' + i))
		options[i] = gen.Option{Key: k, Text: o.text, MisconceptionID: o.mis}
		if o.mis == "" {
			answerKey = k
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
	job := contexts[p.context%len(contexts)].job
	var b strings.Builder
	for i, t := range p.times {
		switch i {
		case 0:
			fmt.Fprintf(&b, "%s dapat %s sendirian dalam %d hari", names[i], job, t)
		case len(p.times) - 1:
			fmt.Fprintf(&b, ", dan %s dalam %d hari", names[i], t)
		default:
			fmt.Fprintf(&b, ", %s dalam %d hari", names[i], t)
		}
	}
	b.WriteString(". ")
	if p.together == 0 {
		b.WriteString("Jika semuanya bekerja bersama, pekerjaan itu selesai dalam ... hari.")
	} else {
		fmt.Fprintf(&b, "A dan B bekerja bersama selama %d hari, kemudian B berhenti dan A melanjutkan sendirian sampai selesai. Seluruh pekerjaan itu selesai dalam ... hari.", p.together)
	}
	return b.String()
}

func daysText(n int) string { return fmt.Sprintf("%d hari", n) }

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int { return a / gcd(a, b) * b }
