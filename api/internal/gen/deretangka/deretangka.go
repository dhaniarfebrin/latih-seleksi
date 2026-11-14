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
)

// Misconception labels attached to wrong options.
const (
	MisBedaTerakhir        = "beda-terakhir"
	MisRasioJadiBeda       = "rasio-jadi-beda"
	MisAbaikanSelangSeling = "abaikan-selang-seling"
	MisSalahHitungDekat    = "salah-hitung-dekat"
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
	kindOrder3               // differences of differences step by a growing amount
)

// params are the numbers behind one question. n terms are shown and the
// (n+1)th is asked. The meaning of the other fields depends on kind:
// a is the first term (fib and alt: the first odd term), b the first even
// term (alt) or second term (fib), d the first difference (alt: the odd
// subsequence step), d2 the even subsequence step, r the ratio, e the
// increment of the differences (order3: of the second differences), s the
// base of the powers and c the constant added to a power (order3: the first
// second difference).
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
	case kindGeom:
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
	case kindStepped, kindOrder3:
		if p.e < 1 {
			return false
		}
	}
	_, ok := seq(p)
	return ok
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
		if rng.Intn(2) == 0 {
			n := rng.Intn(3) + 5
			d := signedNonZero(rng, 8)
			if d > 0 {
				d++
			} else {
				d--
			}
			a := rng.Intn(30) + 1
			if d < 0 {
				a = -d*n + rng.Intn(10) + 1 // keep the asked term positive
			}
			return params{kind: kindArith, n: n, a: a, d: d}
		}
		return params{kind: kindGeom, n: rng.Intn(2) + 5, a: rng.Intn(6) + 1, r: rng.Intn(2) + 2}
	case gen.Band2:
		if rng.Intn(2) == 0 {
			return params{kind: kindStepped, n: rng.Intn(3) + 5, a: rng.Intn(20) + 1, d: rng.Intn(8) + 1, e: rng.Intn(4) + 1}
		}
		return params{kind: kindAlt, n: rng.Intn(2) + 6, a: rng.Intn(30) + 1, d: signedNonZero(rng, 9), b: rng.Intn(30) + 1, d2: signedNonZero(rng, 9)}
	default:
		switch rng.Intn(5) {
		case 0:
			return params{kind: kindSquare, n: rng.Intn(3) + 5, s: rng.Intn(8) + 1, c: rng.Intn(17) - 4}
		case 1:
			return params{kind: kindCube, n: rng.Intn(2) + 5, s: rng.Intn(7) + 1, c: rng.Intn(21) - 10}
		case 2:
			return params{kind: kindFib, n: rng.Intn(2) + 6, a: rng.Intn(9) + 1, b: rng.Intn(9) + 1}
		case 3:
			return params{kind: kindDiffGeom, n: rng.Intn(2) + 5, a: rng.Intn(10) + 1, d: rng.Intn(4) + 1, r: rng.Intn(2) + 2}
		default:
			return params{kind: kindOrder3, n: rng.Intn(2) + 6, a: rng.Intn(20) + 1, d: rng.Intn(6) + 1, c: rng.Intn(4), e: rng.Intn(3) + 1}
		}
	}
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
	case kindOrder3:
		v, x, y := p.a, p.d, p.c
		for k := range t {
			t[k] = v
			v += x
			x += y
			y += p.e
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
	default: // kindOrder3
		ds := diffs(t[:n])
		ss := diffs(ds)
		ns := ss[len(ss)-1] + p.e
		nd := ds[len(ds)-1] + ns
		return []string{
			fmt.Sprintf("Selisih antarsuku: %s.", joined(ds)),
			fmt.Sprintf("Selisih dari selisih itu: %s, bertambah %d setiap langkah.", joined(ss), p.e),
			fmt.Sprintf("Selisih berikutnya = %s.", plus(ds[len(ds)-1], ns)),
			fmt.Sprintf("Suku berikutnya = %s.", plus(last, nd)),
		}
	}
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
