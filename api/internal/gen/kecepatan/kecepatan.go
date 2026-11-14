// Package kecepatan generates relative-speed TIU questions
// (pattern TIU.numerik.kecepatan-relatif).
package kecepatan

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.numerik.kecepatan-relatif"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
	maxHours    = 40 // upper bound of any option, as a sanity limit
)

// Misconception labels attached to wrong options.
const (
	MisJumlahSelisihTertukar = "jumlah-selisih-tertukar"
	MisKecepatanTunggal      = "kecepatan-tunggal"
	MisRataRataKecepatan     = "rata-rata-kecepatan"
	MisLupaSelang            = "lupa-selang-keberangkatan"
	MisSalahHitungDekat      = "salah-hitung-dekat"
)

type kind int

const (
	opposite  kind = iota // leave together, toward each other
	catchUp               // same direction, B ahead of A
	lateStart             // toward each other, B leaves delay hours after A
)

// params are the numbers behind one question: distance in km, speeds of A and
// B in km/jam, and delay in hours (only for lateStart).
type params struct {
	kind     kind
	distance int
	vA, vB   int
	delay    int
	context  int
}

type contextTemplate struct {
	vehicle string
}

var contexts = []contextTemplate{
	{"mobil"},
	{"kereta"},
	{"bus"},
}

// Family implements gen.Family for relative-speed questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Kecepatan relatif" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisJumlahSelisihTertukar: "Menukar aturan: kecepatan dijumlahkan saat searah atau dikurangkan saat berhadapan",
		MisKecepatanTunggal:      "Hanya memakai kecepatan satu kendaraan, mengabaikan gerak kendaraan lainnya",
		MisRataRataKecepatan:     "Memakai rata-rata kedua kecepatan, bukan jumlah atau selisihnya",
		MisLupaSelang:            "Lupa menambahkan selang keberangkatan: waktu dihitung dari saat kendaraan kedua berangkat",
		MisSalahHitungDekat:      "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Parameters whose meeting
// time is not a whole number of hours are rejected; attempt k uses a seed
// derived from (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, wholeHours)
}

func wholeHours(p params) bool {
	h, ok := solveHours(p)
	return ok && h > 0 && h <= maxHours
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("kecepatan: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("kecepatan: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

func speed(rng *rand.Rand) int { return 30 + 5*rng.Intn(11) } // 30..80

// pickParams builds the distance from the intended whole-hour answer, so the
// result is whole by construction.
func pickParams(rng *rand.Rand, band gen.Band) params {
	p := params{vA: speed(rng), vB: speed(rng), context: rng.Intn(len(contexts))}
	t := rng.Intn(4) + 2 // 2..5 hours of closing in
	switch band {
	case gen.Band1:
		p.kind = opposite
		p.distance = t * (p.vA + p.vB)
	case gen.Band2:
		p.kind = catchUp
		if p.vA == p.vB {
			p.vA += 5
		}
		if p.vA < p.vB {
			p.vA, p.vB = p.vB, p.vA
		}
		p.distance = t * (p.vA - p.vB)
	default:
		p.kind = lateStart
		p.delay = rng.Intn(3) + 1 // 1..3 hours head start for A
		p.distance = p.vA*p.delay + t*(p.vA+p.vB)
	}
	return p
}

// solveHours returns the hours until the meeting, counted from A's departure.
// ok is false when the result is not a positive whole number.
func solveHours(p params) (int, bool) {
	switch p.kind {
	case opposite:
		return divide(p.distance, p.vA+p.vB)
	case catchUp:
		return divide(p.distance, p.vA-p.vB)
	case lateStart:
		rest := p.distance - p.vA*p.delay
		h, ok := divide(rest, p.vA+p.vB)
		return h + p.delay, ok
	}
	return 0, false
}

// divide returns num/den when it is a positive whole number.
func divide(num, den int) (int, bool) {
	if den <= 0 || num <= 0 || num%den != 0 {
		return 0, false
	}
	return num / den, true
}

func steps(p params, hours int) []string {
	switch p.kind {
	case opposite:
		return []string{
			fmt.Sprintf("Berhadapan: jarak antar-kendaraan mengecil dengan kecepatan gabungan %d + %d = %d km/jam.", p.vA, p.vB, p.vA+p.vB),
			fmt.Sprintf("Waktu bertemu = %d : %d = %d jam.", p.distance, p.vA+p.vB, hours),
		}
	case catchUp:
		return []string{
			fmt.Sprintf("Searah: A mendekati B dengan kecepatan relatif %d - %d = %d km/jam.", p.vA, p.vB, p.vA-p.vB),
			fmt.Sprintf("Waktu menyusul = %d : %d = %d jam.", p.distance, p.vA-p.vB, hours),
		}
	default:
		head := p.vA * p.delay
		rest := p.distance - head
		return []string{
			fmt.Sprintf("Selama %d jam sebelum B berangkat, A menempuh %d x %d = %d km, sisa jarak %d - %d = %d km.", p.delay, p.vA, p.delay, head, p.distance, head, rest),
			fmt.Sprintf("Setelah B berangkat, kecepatan gabungan %d + %d = %d km/jam, sehingga mereka bertemu %d : %d = %d jam kemudian.", p.vA, p.vB, p.vA+p.vB, rest, p.vA+p.vB, hours-p.delay),
			fmt.Sprintf("Waktu sejak A berangkat = %d + %d = %d jam.", p.delay, hours-p.delay, hours),
		}
	}
}

func build(p params, seed int64) gen.Draft {
	key, _ := solveHours(p)

	type cand struct {
		hours int
		mis   string
	}
	var cands []cand
	if p.kind == lateStart {
		cands = append(cands, cand{key - p.delay, MisLupaSelang})
	}
	swapped := p.vA - p.vB // |difference| when the real rule is the sum
	if swapped < 0 {
		swapped = -swapped
	}
	if p.kind == catchUp {
		swapped = p.vA + p.vB
	}
	cands = append(cands,
		cand{wholeOrZero(p.distance, swapped), MisJumlahSelisihTertukar},
		cand{wholeOrZero(p.distance, p.vA), MisKecepatanTunggal},
		cand{wholeOrZero(p.distance, p.vB), MisKecepatanTunggal},
		cand{wholeOrZero(2*p.distance, p.vA+p.vB), MisRataRataKecepatan},
	)
	for d := 1; d <= maxHours; d++ {
		cands = append(cands, cand{key + d, MisSalahHitungDekat}, cand{key - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	opts := []opt{{hoursText(key), ""}}
	seen := map[string]bool{hoursText(key): true}
	usedMis := map[string]bool{}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		if c.hours <= 0 {
			continue
		}
		t := hoursText(c.hours)
		// one option per specific misconception; salah-hitung-dekat may repeat
		if seen[t] || (c.mis != MisSalahHitungDekat && usedMis[c.mis]) {
			continue
		}
		seen[t], usedMis[c.mis] = true, true
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
		Steps:     steps(p, key),
		Seed:      seed,
	}
}

// wholeOrZero returns num/den when positive and whole, else 0 (skipped).
func wholeOrZero(num, den int) int {
	h, _ := divide(num, den)
	return h
}

func stem(p params) string {
	v := contexts[p.context%len(contexts)].vehicle
	var b strings.Builder
	switch p.kind {
	case opposite:
		fmt.Fprintf(&b, "Kota P dan kota Q berjarak %d km. %s A berangkat dari P menuju Q dengan kecepatan %d km/jam, sedangkan %s B berangkat dari Q menuju P dengan kecepatan %d km/jam. Keduanya berangkat pada waktu yang sama. Mereka bertemu setelah ... jam.",
			p.distance, capitalize(v), p.vA, v, p.vB)
	case catchUp:
		fmt.Fprintf(&b, "%s A melaju dengan kecepatan %d km/jam mengejar %s B yang berada %d km di depannya dan melaju searah dengan kecepatan %d km/jam. A menyusul B setelah ... jam.",
			capitalize(v), p.vA, v, p.distance, p.vB)
	default:
		fmt.Fprintf(&b, "Kota P dan kota Q berjarak %d km. %s A berangkat dari P menuju Q dengan kecepatan %d km/jam. %d jam setelah A berangkat, %s B berangkat dari Q menuju P dengan kecepatan %d km/jam. Mereka bertemu ... jam setelah A berangkat.",
			p.distance, capitalize(v), p.vA, p.delay, v, p.vB)
	}
	return b.String()
}

func hoursText(n int) string { return fmt.Sprintf("%d jam", n) }

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
