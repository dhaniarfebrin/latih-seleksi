// Package ratarata generates combined-average TIU questions
// (pattern TIU.numerik.rata-rata-gabungan).
package ratarata

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.numerik.rata-rata-gabungan"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
	maxTotal    = 60 // upper bound of the head count behind any average
)

// Misconception labels attached to wrong options.
const (
	MisRataRataDuaRata  = "rata-rata-dua-rata"
	MisPembagiSalah     = "pembagi-salah"
	MisSelisihDibalik   = "selisih-dibalik"
	MisSalahHitungDekat = "salah-hitung-dekat"
)

// kind selects the question variant.
type kind int

const (
	kindGroups kind = iota // average of several groups combined
	kindAdd                // one value joins a set with a known average
	kindRemove             // one value leaves a set with a known average
	kindTarget             // value to add so the average becomes target
)

type group struct{ n, r int }

// params are the numbers behind one question. For kindGroups only groups is
// used; for the other kinds n values average r, x is the value that joins or
// leaves, and target is the wanted average (kindTarget only).
type params struct {
	kind    kind
	groups  []group
	n, r, x int
	target  int
	context int
}

type contextTemplate struct {
	measure string // what is averaged
	member  string // what the n counts
	unit    string // suffix of every value, including the leading space
	lo, hi  int    // range of every value shown in the stem
}

var contexts = []contextTemplate{
	{"nilai", "siswa", "", 45, 98},
	{"berat badan", "orang", " kg", 40, 85},
	{"gaji", "karyawan", " juta rupiah", 4, 35},
}

var ordinals = []string{"pertama", "kedua", "ketiga"}

// Family implements gen.Family for combined-average questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Rata-rata gabungan" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisRataRataDuaRata:  "Merata-ratakan rata-rata begitu saja, tanpa memperhitungkan banyak data tiap kelompok",
		MisPembagiSalah:     "Membagi jumlah dengan pembagi yang salah, bukan dengan seluruh banyak data setelah perubahan",
		MisSelisihDibalik:   "Membalik arah perubahan: menambah padahal dikurangi atau sebaliknya",
		MisSalahHitungDekat: "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Parameters whose result is
// not a whole number are rejected; attempt k uses a seed derived from
// (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, valid)
}

func valid(p params) bool {
	ctx := contexts[p.context%len(contexts)]
	inRange := func(v int) bool { return v >= ctx.lo && v <= ctx.hi }
	switch p.kind {
	case kindGroups:
		heads := 0
		for _, g := range p.groups {
			if g.n < 3 || !inRange(g.r) {
				return false
			}
			heads += g.n
		}
		if heads > maxTotal {
			return false
		}
	case kindAdd, kindRemove:
		if !inRange(p.r) || !inRange(p.x) {
			return false
		}
	case kindTarget:
		if !inRange(p.r) || p.target <= 0 || p.target == p.r || !inRange(p.target) {
			return false
		}
	}
	res, ok := solve(p)
	if !ok {
		return false
	}
	if p.kind == kindTarget {
		return inRange(res.key)
	}
	return true
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("ratarata: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("ratarata: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

// signedNonZero returns a value in [-max, max] other than 0.
func signedNonZero(rng *rand.Rand, max int) int {
	v := rng.Intn(max) + 1
	if rng.Intn(2) == 0 {
		return -v
	}
	return v
}

func pickParams(rng *rand.Rand, band gen.Band) params {
	ctx := rng.Intn(len(contexts))
	c := contexts[ctx]
	switch band {
	case gen.Band1:
		// n1*(r1-a) = n2*(a-r2) with n1=m*u, n2=m*v, u and v coprime: choose
		// the combined average a, then r1 = a + v*t and r2 = a - u*t.
		u, v, m := rng.Intn(5)+1, rng.Intn(5)+1, rng.Intn(4)+2
		if gcd(u, v) != 1 {
			return params{context: ctx} // rejected by valid
		}
		a := c.lo + rng.Intn(c.hi-c.lo+1)
		t := signedNonZero(rng, 3)
		return params{kind: kindGroups, context: ctx, groups: []group{{m * u, a + v*t}, {m * v, a - u*t}}}
	case gen.Band2:
		n := rng.Intn(7) + 4
		r := c.lo + rng.Intn(c.hi-c.lo+1)
		e := signedNonZero(rng, 3) // change of the average
		switch rng.Intn(3) {
		case 0:
			return params{kind: kindAdd, context: ctx, n: n, r: r, x: r + e*(n+1)}
		case 1:
			return params{kind: kindRemove, context: ctx, n: n, r: r, x: r - e*(n-1)}
		default:
			return params{kind: kindTarget, context: ctx, n: n, r: r, target: r + e}
		}
	default:
		// deviations d1 = k*t1, d2 = k*t2 from the combined average a, with
		// k the size of group three, make d3 = -(n1*t1 + n2*t2) whole.
		n1, n2, n3 := rng.Intn(7)+3, rng.Intn(7)+3, rng.Intn(4)+3
		t1, t2 := rng.Intn(3)-1, rng.Intn(3)-1
		if t1 == 0 && t2 == 0 {
			return params{context: ctx}
		}
		a := c.lo + rng.Intn(c.hi-c.lo+1)
		return params{kind: kindGroups, context: ctx, groups: []group{
			{n1, a + n3*t1}, {n2, a + n3*t2}, {n3, a - (n1*t1 + n2*t2)},
		}}
	}
}

type result struct {
	key   int
	steps []string
}

// solve returns the key and solution steps; ok is false when params are
// malformed or the key is not a positive whole number.
func solve(p params) (result, bool) {
	ctx := contexts[p.context%len(contexts)]
	val := func(v int) string { return fmt.Sprintf("%d%s", v, ctx.unit) }
	switch p.kind {
	case kindGroups:
		if len(p.groups) < 2 {
			return result{}, false
		}
		var steps, sums, counts []string
		total, heads := 0, 0
		for i, g := range p.groups {
			if g.n <= 0 || g.r <= 0 {
				return result{}, false
			}
			steps = append(steps, fmt.Sprintf("Jumlah kelompok %s = %d x %d = %d.", ordinals[i], g.n, g.r, g.n*g.r))
			sums = append(sums, fmt.Sprint(g.n*g.r))
			counts = append(counts, fmt.Sprint(g.n))
			total += g.n * g.r
			heads += g.n
		}
		if total%heads != 0 {
			return result{}, false
		}
		key := total / heads
		steps = append(steps,
			fmt.Sprintf("Jumlah seluruhnya = %s = %d.", strings.Join(sums, " + "), total),
			fmt.Sprintf("Banyak seluruhnya = %s = %d.", strings.Join(counts, " + "), heads),
			fmt.Sprintf("Rata-rata gabungan = %d : %d = %s.", total, heads, val(key)))
		return result{key, steps}, key > 0
	case kindAdd, kindRemove:
		if p.n < 2 || p.r <= 0 || p.x <= 0 {
			return result{}, false
		}
		sum := p.n * p.r
		steps := []string{fmt.Sprintf("Jumlah semula = %d x %d = %d.", p.n, p.r, sum)}
		var newSum, newN int
		if p.kind == kindAdd {
			newSum, newN = sum+p.x, p.n+1
			steps = append(steps, fmt.Sprintf("Jumlah baru = %d + %d = %d, banyak data = %d + 1 = %d.", sum, p.x, newSum, p.n, newN))
		} else {
			newSum, newN = sum-p.x, p.n-1
			steps = append(steps, fmt.Sprintf("Jumlah baru = %d - %d = %d, banyak data = %d - 1 = %d.", sum, p.x, newSum, p.n, newN))
		}
		if newSum <= 0 || newSum%newN != 0 {
			return result{}, false
		}
		key := newSum / newN
		steps = append(steps, fmt.Sprintf("Rata-rata baru = %d : %d = %s.", newSum, newN, val(key)))
		return result{key, steps}, true
	case kindTarget:
		if p.n < 2 || p.r <= 0 || p.target <= 0 {
			return result{}, false
		}
		sum := p.n * p.r
		want := p.target * (p.n + 1)
		key := want - sum
		steps := []string{
			fmt.Sprintf("Jumlah semula = %d x %d = %d.", p.n, p.r, sum),
			fmt.Sprintf("Jumlah yang diinginkan = %d x (%d + 1) = %d.", p.target, p.n, want),
			fmt.Sprintf("Nilai yang ditambahkan = %d - %d = %s.", want, sum, val(key)),
		}
		return result{key, steps}, key > 0
	}
	return result{}, false
}

func build(p params, seed int64) gen.Draft {
	ctx := contexts[p.context%len(contexts)]
	res, _ := solve(p)
	key := res.key

	type cand struct {
		v   int
		mis string
	}
	var cands []cand
	// whole is a candidate only when num divides evenly and stays positive.
	whole := func(num, den int, mis string) {
		if den > 0 && num > 0 && num%den == 0 {
			cands = append(cands, cand{num / den, mis})
		}
	}
	switch p.kind {
	case kindGroups:
		total, avgs := 0, 0
		for _, g := range p.groups {
			total += g.n * g.r
			avgs += g.r
		}
		whole(avgs, len(p.groups), MisRataRataDuaRata)
		whole(total, len(p.groups), MisPembagiSalah)
	case kindAdd:
		sum := p.n * p.r
		whole(p.r+p.x, 2, MisRataRataDuaRata)
		whole(sum+p.x, p.n, MisPembagiSalah)
		whole(sum-p.x, p.n-1, MisSelisihDibalik)
	case kindRemove:
		sum := p.n * p.r
		whole(sum-p.x, p.n, MisPembagiSalah)
		whole(sum+p.x, p.n+1, MisSelisihDibalik)
	case kindTarget:
		sum := p.n * p.r
		whole(2*p.target-p.r, 1, MisRataRataDuaRata)
		whole(p.n*p.target-sum, 1, MisPembagiSalah)
		whole(p.r*(p.n+1)-p.n*p.target, 1, MisSelisihDibalik)
	}
	for d := 1; d <= key+optionCount; d++ {
		cands = append(cands, cand{key + d, MisSalahHitungDekat}, cand{key - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	text := func(v int) string { return fmt.Sprintf("%d%s", v, ctx.unit) }
	opts := []opt{{text(key), ""}}
	seen := map[string]bool{text(key): true}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		if c.v <= 0 {
			continue
		}
		t := text(c.v)
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
		Steps:     res.steps,
		Seed:      seed,
	}
}

func stem(p params) string {
	ctx := contexts[p.context%len(contexts)]
	var b strings.Builder
	switch p.kind {
	case kindGroups:
		for i, g := range p.groups {
			fmt.Fprintf(&b, "Kelompok %s terdiri dari %d %s dengan rata-rata %s %d%s. ", ordinals[i], g.n, ctx.member, ctx.measure, g.r, ctx.unit)
		}
		fmt.Fprintf(&b, "Rata-rata %s gabungan seluruh kelompok itu adalah ...", ctx.measure)
	case kindAdd:
		fmt.Fprintf(&b, "Rata-rata %s dari %d %s adalah %d%s. ", ctx.measure, p.n, ctx.member, p.r, ctx.unit)
		fmt.Fprintf(&b, "Jika satu %s lagi dengan %s %d%s ikut bergabung, rata-rata %s yang baru adalah ...", ctx.member, ctx.measure, p.x, ctx.unit, ctx.measure)
	case kindRemove:
		fmt.Fprintf(&b, "Rata-rata %s dari %d %s adalah %d%s. ", ctx.measure, p.n, ctx.member, p.r, ctx.unit)
		fmt.Fprintf(&b, "Jika satu %s dengan %s %d%s keluar, rata-rata %s yang baru adalah ...", ctx.member, ctx.measure, p.x, ctx.unit, ctx.measure)
	case kindTarget:
		fmt.Fprintf(&b, "Rata-rata %s dari %d %s adalah %d%s. ", ctx.measure, p.n, ctx.member, p.r, ctx.unit)
		fmt.Fprintf(&b, "Satu %s lagi akan ditambahkan agar rata-rata %s menjadi %d%s. Besar %s yang ditambahkan itu adalah ...", ctx.member, ctx.measure, p.target, ctx.unit, ctx.measure)
	}
	return b.String()
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
