// Package deretangka generates number-series TIU questions
// (pattern TIU.numerik.deret-angka).
package deretangka

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.numerik.deret-angka"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
	maxTerm     = 2000 // every term, shown or asked, stays within [1, maxTerm]
	maxPrime    = 200  // prime series stay small, as in a real exam
)

// Misconception labels attached to wrong options.
const (
	MisBedaTerakhir        = "beda-terakhir"
	MisRasioJadiBeda       = "rasio-jadi-beda"
	MisAbaikanSelangSeling = "abaikan-selang-seling"
	MisSalahHitungDekat    = "salah-hitung-dekat"
	MisSatuOperasi         = "satu-operasi"
	MisBukanPrima          = "bukan-prima"
)

// kind selects the rule behind the series.
type kind int

const (
	kindArith    kind = iota // constant difference d
	kindGeom                 // constant ratio r
	kindStepped              // differences d, d+e, d+2e, ...
	kindAlt                  // odd and even terms follow separate arithmetic rules
	kindSquare               // (s+k)^2 + c
	kindCube                 // (s+k)^3 + c
	kindFib                  // each term is the sum of the two before it
	kindDiffGeom             // differences d, d*r, d*r^2, ...
	kindGeomDown             // constant ratio 1/r: each term is divided by r
	kindAltOps               // two operations (add/subtract d, multiply by r) applied in turn
	kindPrime                // consecutive prime numbers
)

// params are the numbers behind one question. n terms are shown and the
// (n+1)th is asked. The meaning of the other fields depends on kind:
// a is the first term (fib and alt: the first odd term), b the first even
// term (alt) or second term (fib), d the first difference (alt: the odd
// subsequence step; altops: the amount added or subtracted), d2 the even
// subsequence step, r the ratio (geomdown: the divisor; altops: the
// multiplier), e the increment of the differences (altops: the variant, see
// altOp), s the base of the powers and c the constant added to a power.
// For prime, a is the first prime of the series.
type params struct {
	kind                    kind
	n                       int
	a, b, d, d2, r, e, s, c int
}

// Family implements gen.Family for number-series questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Deret angka" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisBedaTerakhir:        "Mengulang beda terakhir saja, padahal beda antarsuku ikut berubah mengikuti pola",
		MisRasioJadiBeda:       "Tertukar antara menambah dan mengalikan: rasio dipakai sebagai beda, atau beda dipakai sebagai rasio",
		MisAbaikanSelangSeling: "Menganggap semua suku mengikuti satu aturan, padahal suku ganjil dan suku genap punya aturan sendiri",
		MisSalahHitungDekat:    "Salah hitung kecil: hasil mendekati jawaban benar",
		MisSatuOperasi:         "Hanya satu dari dua operasi yang diulang, padahal kedua operasi dipakai bergantian",
		MisBukanPrima:          "Memilih bilangan berikutnya yang bukan prima, padahal harus bilangan prima setelah suku terakhir",
	}
}

// Generate returns the question for seed and band. Parameters whose terms
// leave [1, maxTerm] are rejected; attempt k uses a seed derived from
// (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, valid)
}

func valid(p params) bool {
	switch p.kind {
	case kindArith:
		if p.d == 0 {
			return false
		}
	case kindGeom, kindGeomDown:
		if p.r < 2 || p.r > 3 {
			return false
		}
	case kindDiffGeom:
		// a*(r-1) == d makes the series itself geometric, a band 1 rule
		if p.r < 2 || p.r > 3 || p.a*(p.r-1) == p.d {
			return false
		}
	case kindAlt:
		if p.d == 0 || p.d2 == 0 || p.d == p.d2 {
			return false
		}
	case kindStepped:
		if p.e < 1 {
			return false
		}
	case kindAltOps:
		if p.d < 1 || p.r < 2 || p.r > 3 || p.e < 0 || p.e > 3 {
			return false
		}
	}
	t, ok := seq(p)
	if !ok {
		return false
	}
	switch p.kind {
	case kindAltOps:
		// a series that repeats itself every two terms is not a two-operation rule
		return t[0] != t[2]
	case kindPrime:
		return !fitsSimpleRule(t[:p.n])
	}
	return true
}

// fitsSimpleRule reports whether the shown terms also fit a constant
// difference of order 1-3 or two arithmetic subsequences, which would make
// a prime series ambiguous.
func fitsSimpleRule(shown []int) bool {
	level := shown
	for k := 1; k <= 3; k++ {
		level = diffs(level)
		if constant(level) {
			return true
		}
	}
	var odd, even []int
	for i, v := range shown {
		if i%2 == 0 {
			odd = append(odd, v)
		} else {
			even = append(even, v)
		}
	}
	return constant(diffs(odd)) && constant(diffs(even))
}

// constant reports whether xs has at least two entries, all equal.
func constant(xs []int) bool {
	for _, x := range xs {
		if x != xs[0] {
			return false
		}
	}
	return len(xs) >= 2
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("deretangka: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("deretangka: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
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
	switch band {
	case gen.Band1:
		switch rng.Intn(4) {
		case 0, 1:
			n := rng.Intn(3) + 5
			d := signedNonZero(rng, 14)
			if d > 0 {
				d++
			} else {
				d--
			}
			a := rng.Intn(60) + 1
			if d < 0 {
				a = -d*n + rng.Intn(20) + 1 // keep the asked term positive
			}
			return params{kind: kindArith, n: n, a: a, d: d}
		case 2:
			return params{kind: kindGeom, n: rng.Intn(2) + 5, a: rng.Intn(12) + 1, r: rng.Intn(2) + 2}
		default:
			n := rng.Intn(2) + 5
			r := rng.Intn(2) + 2
			pow := 1
			for i := 0; i < n; i++ {
				pow *= r
			}
			// the first term is a multiple of r^n so every term, the asked one included, is whole
			return params{kind: kindGeomDown, n: n, a: (rng.Intn(maxTerm/pow) + 1) * pow, r: r}
		}
	case gen.Band2:
		switch rng.Intn(4) {
		case 0:
			return params{kind: kindStepped, n: rng.Intn(3) + 5, a: rng.Intn(20) + 1, d: rng.Intn(8) + 1, e: rng.Intn(4) + 1}
		case 1:
			return params{kind: kindAlt, n: rng.Intn(2) + 6, a: rng.Intn(30) + 1, d: signedNonZero(rng, 9), b: rng.Intn(30) + 1, d2: signedNonZero(rng, 9)}
		case 2:
			return params{kind: kindAltOps, n: rng.Intn(2) + 6, a: rng.Intn(20) + 1, d: rng.Intn(9) + 1, r: rng.Intn(2) + 2, e: rng.Intn(4)}
		default:
			return pickPrime(rng)
		}
	default:
		switch rng.Intn(5) {
		case 0:
			return params{kind: kindSquare, n: rng.Intn(3) + 5, s: rng.Intn(30) + 1, c: rng.Intn(51) - 10}
		case 1:
			return params{kind: kindCube, n: rng.Intn(2) + 5, s: rng.Intn(8) + 1, c: rng.Intn(41) - 20}
		case 2:
			return params{kind: kindFib, n: rng.Intn(2) + 6, a: rng.Intn(30) + 1, b: rng.Intn(30) + 1}
		case 3:
			return params{kind: kindDiffGeom, n: rng.Intn(2) + 5, a: rng.Intn(15) + 1, d: rng.Intn(9) + 1, r: rng.Intn(2) + 2}
		default:
			return pickPrime(rng)
		}
	}
}

// primes lists every prime in [2, maxPrime].
var primes = func() []int {
	var out []int
	for v := 2; v <= maxPrime; v++ {
		if isPrime(v) {
			out = append(out, v)
		}
	}
	return out
}()

func pickPrime(rng *rand.Rand) params {
	n := rng.Intn(3) + 5
	return params{kind: kindPrime, n: n, a: primes[rng.Intn(len(primes)-n)]}
}

func isPrime(v int) bool {
	if v < 2 {
		return false
	}
	for d := 2; d*d <= v; d++ {
		if v%d == 0 {
			return false
		}
	}
	return true
}

func nextPrime(v int) int {
	for v++; !isPrime(v); v++ {
	}
	return v
}

// altOp returns the operation that turns term k into term k+1 for altops:
// when mul is set amt is the multiplier, otherwise it is added to the term.
// Variant p.e is 0: +d then x r; 1: x r then +d; 2: x r then -d; 3: -d then x r.
func altOp(p params, k int) (mul bool, amt int) {
	first := k%2 == 0
	switch p.e {
	case 0:
		if first {
			return false, p.d
		}
	case 1:
		if first {
			return true, p.r
		}
		return false, p.d
	case 2:
		if first {
			return true, p.r
		}
		return false, -p.d
	default:
		if first {
			return false, -p.d
		}
	}
	return true, p.r
}

func applyOp(v int, mul bool, amt int) int {
	if mul {
		return v * amt
	}
	return v + amt
}

// seq returns the n+1 terms of the series (the last one is the key); ok is
// false when the kind is unknown or any term leaves [1, maxTerm].
func seq(p params) ([]int, bool) {
	if p.n < 5 {
		return nil, false
	}
	t := make([]int, p.n+1)
	switch p.kind {
	case kindArith:
		for k := range t {
			t[k] = p.a + k*p.d
		}
	case kindGeom:
		v := p.a
		for k := range t {
			t[k] = v
			v *= p.r
		}
	case kindStepped:
		v := p.a
		for k := range t {
			t[k] = v
			v += p.d + k*p.e
		}
	case kindAlt:
		for k := range t {
			if k%2 == 0 {
				t[k] = p.a + (k/2)*p.d
			} else {
				t[k] = p.b + (k/2)*p.d2
			}
		}
	case kindSquare:
		for k := range t {
			m := p.s + k
			t[k] = m*m + p.c
		}
	case kindCube:
		for k := range t {
			m := p.s + k
			t[k] = m*m*m + p.c
		}
	case kindFib:
		t[0], t[1] = p.a, p.b
		for k := 2; k < len(t); k++ {
			t[k] = t[k-1] + t[k-2]
		}
	case kindDiffGeom:
		v, x := p.a, p.d
		for k := range t {
			t[k] = v
			v += x
			x *= p.r
		}
	case kindGeomDown:
		v := p.a
		for k := range t {
			t[k] = v
			if k == len(t)-1 {
				break
			}
			if v%p.r != 0 {
				return nil, false
			}
			v /= p.r
		}
	case kindAltOps:
		t[0] = p.a
		for k := 1; k < len(t); k++ {
			mul, amt := altOp(p, k-1)
			t[k] = applyOp(t[k-1], mul, amt)
		}
	case kindPrime:
		if !isPrime(p.a) {
			return nil, false
		}
		t[0] = p.a
		for k := 1; k < len(t); k++ {
			t[k] = nextPrime(t[k-1])
			if t[k] > maxPrime {
				return nil, false
			}
		}
	default:
		return nil, false
	}
	for _, v := range t {
		if v < 1 || v > maxTerm {
			return nil, false
		}
	}
	return t, true
}

func joined(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = strconv.Itoa(x)
	}
	return strings.Join(s, ", ")
}

// diffs returns the consecutive differences of xs.
func diffs(xs []int) []int {
	out := make([]int, len(xs)-1)
	for i := range out {
		out[i] = xs[i+1] - xs[i]
	}
	return out
}

// plus renders "x + y = x+y" or "x - |y| = x+y" for negative y.
func plus(x, y int) string {
	if y < 0 {
		return fmt.Sprintf("%d - %d = %d", x, -y, x+y)
	}
	return fmt.Sprintf("%d + %d = %d", x, y, x+y)
}

// offset renders the constant added to a power: "", " + 3" or " - 3".
func offset(c int) string {
	switch {
	case c > 0:
		return fmt.Sprintf(" + %d", c)
	case c < 0:
		return fmt.Sprintf(" - %d", -c)
	}
	return ""
}

func steps(p params, t []int) []string {
	n := p.n
	last := t[n-1]
	key := t[n]
	switch p.kind {
	case kindArith:
		return []string{
			fmt.Sprintf("Selisih antarsuku selalu %d.", p.d),
			fmt.Sprintf("Suku berikutnya = %s.", plus(last, p.d)),
		}
	case kindGeom:
		return []string{
			fmt.Sprintf("Setiap suku dikalikan %d untuk mendapat suku berikutnya.", p.r),
			fmt.Sprintf("Suku berikutnya = %d x %d = %d.", last, p.r, key),
		}
	case kindStepped:
		ds := diffs(t[:n])
		next := ds[len(ds)-1] + p.e
		return []string{
			fmt.Sprintf("Selisih antarsuku: %s, bertambah %d setiap langkah.", joined(ds), p.e),
			fmt.Sprintf("Selisih berikutnya = %s.", plus(ds[len(ds)-1], p.e)),
			fmt.Sprintf("Suku berikutnya = %s.", plus(last, next)),
		}
	case kindAlt:
		var odd, even []int
		for k := 0; k < n; k++ {
			if k%2 == 0 {
				odd = append(odd, t[k])
			} else {
				even = append(even, t[k])
			}
		}
		own, step, which := odd, p.d, "ganjil"
		if n%2 == 1 {
			own, step, which = even, p.d2, "genap"
		}
		return []string{
			fmt.Sprintf("Suku ke-1, ke-3, dan seterusnya: %s, selisihnya %d.", joined(odd), p.d),
			fmt.Sprintf("Suku ke-2, ke-4, dan seterusnya: %s, selisihnya %d.", joined(even), p.d2),
			fmt.Sprintf("Suku ke-%d termasuk deret %s, jadi suku berikutnya = %s.", n+1, which, plus(own[len(own)-1], step)),
		}
	case kindSquare:
		return []string{
			fmt.Sprintf("Suku-suku itu berbentuk n^2%s dengan n = %d, %d, dan seterusnya.", offset(p.c), p.s, p.s+1),
			fmt.Sprintf("Suku berikutnya = %d^2%s = %d.", p.s+n, offset(p.c), key),
		}
	case kindCube:
		return []string{
			fmt.Sprintf("Suku-suku itu berbentuk n^3%s dengan n = %d, %d, dan seterusnya.", offset(p.c), p.s, p.s+1),
			fmt.Sprintf("Suku berikutnya = %d^3%s = %d.", p.s+n, offset(p.c), key),
		}
	case kindFib:
		return []string{
			"Setiap suku adalah jumlah dua suku sebelumnya.",
			fmt.Sprintf("Suku berikutnya = %d + %d = %d.", t[n-2], last, key),
		}
	case kindDiffGeom:
		ds := diffs(t[:n])
		next := ds[len(ds)-1] * p.r
		return []string{
			fmt.Sprintf("Selisih antarsuku: %s, dikalikan %d setiap langkah.", joined(ds), p.r),
			fmt.Sprintf("Selisih berikutnya = %d x %d = %d.", ds[len(ds)-1], p.r, next),
			fmt.Sprintf("Suku berikutnya = %s.", plus(last, next)),
		}
	case kindGeomDown:
		return []string{
			fmt.Sprintf("Setiap suku dibagi %d untuk mendapat suku berikutnya.", p.r),
			fmt.Sprintf("Suku berikutnya = %d : %d = %d.", last, p.r, key),
		}
	case kindAltOps:
		m0, a0 := altOp(p, 0)
		m1, a1 := altOp(p, 1)
		mn, an := altOp(p, n-1)
		return []string{
			fmt.Sprintf("Dua operasi dipakai bergantian: %s, lalu %s, dan seterusnya.", opName(m0, a0), opName(m1, a1)),
			fmt.Sprintf("Contoh dari suku pertama: %s, lalu %s.", opEq(t[0], m0, a0), opEq(t[1], m1, a1)),
			fmt.Sprintf("Operasi berikutnya adalah %s, jadi suku berikutnya = %s.", opName(mn, an), opEq(last, mn, an)),
		}
	default: // kindPrime
		var between []string
		for v := last + 1; v < key; v++ {
			between = append(between, fmt.Sprintf("%d (habis dibagi %d)", v, smallestDivisor(v)))
		}
		gap := fmt.Sprintf("Bilangan setelah %d langsung prima.", last)
		if len(between) > 0 {
			gap = fmt.Sprintf("Bilangan di antara %d dan %d bukan prima: %s.", last, key, strings.Join(between, ", "))
		}
		return []string{
			fmt.Sprintf("Semua suku adalah bilangan prima berurutan: %s.", joined(t[:n])),
			gap,
			fmt.Sprintf("Suku berikutnya adalah bilangan prima setelah %d, yaitu %d.", last, key),
		}
	}
}

// smallestDivisor returns the smallest divisor of v greater than 1.
func smallestDivisor(v int) int {
	d := 2
	for d*d <= v && v%d != 0 {
		d++
	}
	if d*d > v {
		return v
	}
	return d
}

// opName renders an altops operation: "tambah 3", "kurang 3" or "kali 2".
func opName(mul bool, amt int) string {
	switch {
	case mul:
		return fmt.Sprintf("kali %d", amt)
	case amt < 0:
		return fmt.Sprintf("kurang %d", -amt)
	}
	return fmt.Sprintf("tambah %d", amt)
}

// opEq renders applying an altops operation to v as an equation.
func opEq(v int, mul bool, amt int) string {
	if mul {
		return fmt.Sprintf("%d x %d = %d", v, amt, v*amt)
	}
	return plus(v, amt)
}

func build(p params, seed int64) gen.Draft {
	t, _ := seq(p)
	n := p.n
	last, prev := t[n-1], t[n-2]
	key := t[n]

	type cand struct {
		v   int
		mis string
	}
	var cands []cand
	lastDiff := last - prev
	switch p.kind {
	case kindArith:
		if p.d >= 2 {
			cands = append(cands, cand{last * p.d, MisRasioJadiBeda})
		}
	case kindGeom:
		cands = append(cands, cand{last + p.r, MisRasioJadiBeda}, cand{last + lastDiff, MisBedaTerakhir})
	case kindAlt:
		own := p.d
		if n%2 == 1 {
			own = p.d2
		}
		cands = append(cands, cand{last + lastDiff, MisAbaikanSelangSeling}, cand{last + own, MisAbaikanSelangSeling})
	case kindDiffGeom:
		cands = append(cands, cand{last + lastDiff, MisBedaTerakhir}, cand{last + lastDiff + p.r, MisRasioJadiBeda})
	case kindGeomDown:
		cands = append(cands, cand{last - p.r, MisRasioJadiBeda}, cand{last + lastDiff, MisBedaTerakhir})
	case kindAltOps:
		// repeating the previous operation uses only one of the two operations
		mul, amt := altOp(p, n-2)
		cands = append(cands, cand{applyOp(last, mul, amt), MisSatuOperasi}, cand{last + lastDiff, MisBedaTerakhir})
	case kindPrime:
		// the nearest odd numbers above the last term that are not prime
		for v, found := last+1, 0; found < 2; v++ {
			if v%2 == 1 && !isPrime(v) {
				cands = append(cands, cand{v, MisBukanPrima})
				found++
			}
		}
		cands = append(cands, cand{last + lastDiff, MisBedaTerakhir})
	default:
		cands = append(cands, cand{last + lastDiff, MisBedaTerakhir})
	}
	for d := 1; d <= optionCount; d++ {
		cands = append(cands, cand{key + d, MisSalahHitungDekat}, cand{key - d, MisSalahHitungDekat})
	}

	type opt struct{ text, mis string }
	text := func(v int) string { return strconv.Itoa(v) }
	opts := []opt{{text(key), ""}}
	seen := map[string]bool{text(key): true}
	for _, c := range cands {
		if len(opts) == optionCount {
			break
		}
		if c.v <= 0 {
			continue
		}
		s := text(c.v)
		if seen[s] {
			continue
		}
		seen[s] = true
		opts = append(opts, opt{s, c.mis})
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
		Stem:      stem(t[:n]),
		Options:   options,
		AnswerKey: answerKey,
		Steps:     steps(p, t),
		Seed:      seed,
	}
}

func stem(shown []int) string {
	return fmt.Sprintf("Perhatikan deret angka berikut: %s, ... Angka berikutnya dalam deret itu adalah ...", joined(shown))
}
