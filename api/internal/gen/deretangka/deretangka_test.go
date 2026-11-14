package deretangka

import (
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
)

var (
	stemRe   = regexp.MustCompile(`^Perhatikan deret angka berikut: ((?:\d+, )+)\.\.\. `)
	optionRe = regexp.MustCompile(`^(\d+)$`)
)

func atoi(s string) int64 {
	v, _ := strconv.Atoi(s)
	return int64(v)
}

// shownTerms parses the displayed terms from the stem text only.
func shownTerms(t *testing.T, stem string) []int64 {
	t.Helper()
	m := stemRe.FindStringSubmatch(stem)
	if m == nil {
		t.Fatalf("unrecognised stem %q", stem)
	}
	var out []int64
	for _, f := range strings.Split(strings.TrimSuffix(strings.TrimSpace(m[1]), ","), ", ") {
		out = append(out, atoi(strings.TrimSpace(f)))
	}
	return out
}

// prediction is what one rule hypothesis says the next term is.
type prediction struct {
	rule string
	next *big.Rat
}

func rat(v int64) *big.Rat { return big.NewRat(v, 1) }

func diffsOf(xs []int64) []int64 {
	out := make([]int64, len(xs)-1)
	for i := range out {
		out[i] = xs[i+1] - xs[i]
	}
	return out
}

func allEqual(xs []int64) bool {
	for _, x := range xs {
		if x != xs[0] {
			return false
		}
	}
	return len(xs) >= 2
}

// isGeometric reports whether every term squares to the product of its neighbours.
func isGeometric(xs []int64) bool {
	if len(xs) < 3 {
		return false
	}
	for i := 1; i < len(xs)-1; i++ {
		if xs[i] == 0 || xs[i+1]*xs[i-1] != xs[i]*xs[i] {
			return false
		}
	}
	return true
}

// nextGeometric continues a geometric series: last^2 / second-last.
func nextGeometric(xs []int64) *big.Rat {
	l, p := xs[len(xs)-1], xs[len(xs)-2]
	return big.NewRat(l*l, p)
}

// hypotheses lists every rule from the independent detector set that fits
// the shown terms, with the term each would continue with. Order is the
// priority used to name the primary rule of a series.
func hypotheses(ts []int64) []prediction {
	var out []prediction
	n := len(ts)
	// constant k-th difference, k = 1..3 (arithmetic, stepped, cubic, order 3)
	cols := [][]int64{ts}
	for k := 1; k <= 3; k++ {
		level := diffsOf(cols[k-1])
		cols = append(cols, level)
		if !allEqual(level) {
			continue
		}
		// next term = constant k-th difference + last value of every lower level
		v := level[0]
		for j := 0; j < k; j++ {
			v += cols[j][len(cols[j])-1]
		}
		out = append(out, prediction{[]string{"", "arith", "order2", "order3"}[k], rat(v)})
	}
	if isGeometric(ts) {
		out = append(out, prediction{"geom", nextGeometric(ts)})
	}
	// odd/even subsequences both arithmetic
	var odd, even []int64
	for i, v := range ts {
		if i%2 == 0 {
			odd = append(odd, v)
		} else {
			even = append(even, v)
		}
	}
	if allEqual(diffsOf(odd)) && allEqual(diffsOf(even)) {
		own := odd
		if n%2 == 1 {
			own = even
		}
		step := own[len(own)-1] - own[len(own)-2]
		out = append(out, prediction{"alt", rat(own[len(own)-1] + step)})
	}
	// each term is the sum of the two before it
	fib := true
	for i := 2; i < n; i++ {
		if ts[i] != ts[i-1]+ts[i-2] {
			fib = false
		}
	}
	if fib {
		out = append(out, prediction{"fib", rat(ts[n-1] + ts[n-2])})
	}
	// geometric differences
	if ds := diffsOf(ts); len(ds) >= 3 && isGeometric(ds) {
		out = append(out, prediction{"diffgeom", new(big.Rat).Add(rat(ts[n-1]), nextGeometric(ds))})
	}
	return out
}

// primary names the rule of a series: the first fitting hypothesis, using the
// priority arith, geom, alt, order2, fib, diffgeom, order3.
func primary(hs []prediction) string {
	for _, want := range []string{"arith", "geom", "alt", "order2", "fib", "diffgeom", "order3"} {
		for _, h := range hs {
			if h.rule == want {
				return want
			}
		}
	}
	return ""
}

func optionValue(t *testing.T, text string) *big.Rat {
	t.Helper()
	m := optionRe.FindStringSubmatch(text)
	if m == nil {
		t.Fatalf("unparsable option %q", text)
	}
	return rat(atoi(m[1]))
}

func TestPropertyKeyMatchesIndependentRules(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			ts := shownTerms(t, d.Stem)
			if len(ts) < 5 || len(ts) > 7 {
				t.Fatalf("band %d seed %d: %d terms shown", band, seed, len(ts))
			}
			hs := hypotheses(ts)
			if len(hs) == 0 {
				t.Fatalf("band %d seed %d: no independent rule fits %v", band, seed, ts)
			}
			var want *big.Rat
			for _, h := range hs {
				if want == nil {
					want = h.next
				}
				if h.next.Cmp(want) != 0 {
					t.Fatalf("band %d seed %d: ambiguous series %v: %s says %s, %s says %s", band, seed, ts, hs[0].rule, want.RatString(), h.rule, h.next.RatString())
				}
			}
			if !want.IsInt() || want.Sign() <= 0 {
				t.Fatalf("band %d seed %d: next term %s is not a positive whole number", band, seed, want.RatString())
			}
			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			matches := 0
			texts := map[string]bool{}
			for _, o := range d.Options {
				v := optionValue(t, o.Text)
				if texts[o.Text] {
					t.Fatalf("band %d seed %d: duplicate option %q", band, seed, o.Text)
				}
				texts[o.Text] = true
				if v.Cmp(want) == 0 {
					matches++
					if o.Key != d.AnswerKey {
						t.Fatalf("band %d seed %d: AnswerKey %s but value matches %s", band, seed, d.AnswerKey, o.Key)
					}
					if o.MisconceptionID != "" {
						t.Fatalf("band %d seed %d: key option has misconception", band, seed)
					}
				} else if o.MisconceptionID == "" {
					t.Fatalf("band %d seed %d: distractor %s has no misconception", band, seed, o.Key)
				}
			}
			if matches != 1 {
				t.Fatalf("band %d seed %d: %d options equal key %s", band, seed, matches, want.RatString())
			}
			if d.Seed != seed || d.PatternID != PatternID {
				t.Fatalf("band %d seed %d: bad metadata %d %s", band, seed, d.Seed, d.PatternID)
			}
		}
	}
}

func TestBandsEmitExpectedRules(t *testing.T) {
	want := map[gen.Band][]string{
		gen.Band1: {"arith", "geom"},
		gen.Band2: {"order2", "alt"},
		gen.Band3: {"order2", "order3", "fib", "diffgeom"},
	}
	for band, rules := range want {
		seen := map[string]int{}
		for seed := int64(1); seed <= 1000; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			seen[primary(hypotheses(shownTerms(t, d.Stem)))]++
		}
		for _, r := range rules {
			if seen[r] == 0 {
				t.Errorf("band %d never emitted rule %q (seen %v)", band, r, seen)
			}
		}
		if len(seen) != len(rules) {
			t.Errorf("band %d emitted unexpected rules %v", band, seen)
		}
	}
}

func TestDeterminismAndVariation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		stems := map[string]bool{}
		for seed := int64(1); seed <= 200; seed++ {
			a, err1 := fam.Generate(seed, band)
			b, err2 := fam.Generate(seed, band)
			if err1 != nil || err2 != nil {
				t.Fatal(err1, err2)
			}
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("band %d seed %d: drafts differ", band, seed)
			}
			stems[a.Stem] = true
		}
		t.Logf("band %d: %d distinct stems in 200 seeds", band, len(stems))
		if len(stems) < 80 {
			t.Errorf("band %d: only %d distinct stems in 200 seeds", band, len(stems))
		}
	}
}

// byMisconception groups distractor texts by misconception id and returns the key text.
func byMisconception(d gen.Draft) (string, map[string][]string) {
	out := map[string][]string{}
	key := ""
	for _, o := range d.Options {
		if o.Key == d.AnswerKey {
			key = o.Text
			continue
		}
		out[o.MisconceptionID] = append(out[o.MisconceptionID], o.Text)
	}
	for _, v := range out {
		sort.Strings(v)
	}
	return key, out
}

func TestFixtures(t *testing.T) {
	cases := []struct {
		name     string
		p        params
		stem     string
		key      string
		wantMis  map[string][]string
		stepPart string
	}{
		{
			name: "band1 aritmetika naik",
			p:    params{kind: kindArith, n: 5, a: 3, d: 4},
			stem: "Perhatikan deret angka berikut: 3, 7, 11, 15, 19, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "23",
			wantMis: map[string][]string{
				MisRasioJadiBeda:    {"76"},
				MisSalahHitungDekat: {"22", "24", "25"},
			},
			stepPart: "Suku berikutnya = 19 + 4 = 23.",
		},
		{
			name: "band1 aritmetika turun",
			p:    params{kind: kindArith, n: 5, a: 30, d: -4},
			stem: "Perhatikan deret angka berikut: 30, 26, 22, 18, 14, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "10",
			wantMis: map[string][]string{
				MisSalahHitungDekat: {"11", "12", "8", "9"},
			},
			stepPart: "Suku berikutnya = 14 - 4 = 10.",
		},
		{
			name: "band1 geometri",
			p:    params{kind: kindGeom, n: 5, a: 2, r: 3},
			stem: "Perhatikan deret angka berikut: 2, 6, 18, 54, 162, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "486",
			wantMis: map[string][]string{
				MisRasioJadiBeda:    {"165"},
				MisBedaTerakhir:     {"270"},
				MisSalahHitungDekat: {"485", "487"},
			},
			stepPart: "Suku berikutnya = 162 x 3 = 486.",
		},
		{
			name: "band2 beda bertingkat",
			p:    params{kind: kindStepped, n: 5, a: 1, d: 2, e: 3},
			stem: "Perhatikan deret angka berikut: 1, 3, 8, 16, 27, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "41",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"38"},
				MisSalahHitungDekat: {"40", "42", "43"},
			},
			stepPart: "Suku berikutnya = 27 + 14 = 41.",
		},
		{
			name: "band2 selang-seling suku ganjil",
			p:    params{kind: kindAlt, n: 6, a: 5, d: 3, b: 20, d2: -4},
			stem: "Perhatikan deret angka berikut: 5, 20, 8, 16, 11, 12, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "14",
			wantMis: map[string][]string{
				MisAbaikanSelangSeling: {"13", "15"},
				MisSalahHitungDekat:    {"12", "16"},
			},
			stepPart: "Suku ke-7 termasuk deret ganjil, jadi suku berikutnya = 11 + 3 = 14.",
		},
		{
			name: "band2 selang-seling suku genap",
			p:    params{kind: kindAlt, n: 7, a: 22, d: 7, b: 22, d2: 6},
			stem: "Perhatikan deret angka berikut: 22, 22, 29, 28, 36, 34, 43, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "40",
			wantMis: map[string][]string{
				MisAbaikanSelangSeling: {"49", "52"},
				MisSalahHitungDekat:    {"39", "41"},
			},
			stepPart: "Suku ke-8 termasuk deret genap, jadi suku berikutnya = 34 + 6 = 40.",
		},
		{
			name: "band3 kuadrat",
			p:    params{kind: kindSquare, n: 5, s: 2, c: 1},
			stem: "Perhatikan deret angka berikut: 5, 10, 17, 26, 37, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "50",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"48"},
				MisSalahHitungDekat: {"49", "51", "52"},
			},
			stepPart: "Suku berikutnya = 7^2 + 1 = 50.",
		},
		{
			name: "band3 kubik",
			p:    params{kind: kindCube, n: 5, s: 2, c: -1},
			stem: "Perhatikan deret angka berikut: 7, 26, 63, 124, 215, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "342",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"306"},
				MisSalahHitungDekat: {"341", "343", "344"},
			},
			stepPart: "Suku berikutnya = 7^3 - 1 = 342.",
		},
		{
			name: "band3 fibonacci",
			p:    params{kind: kindFib, n: 6, a: 2, b: 3},
			stem: "Perhatikan deret angka berikut: 2, 3, 5, 8, 13, 21, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "34",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"29"},
				MisSalahHitungDekat: {"33", "35", "36"},
			},
			stepPart: "Suku berikutnya = 13 + 21 = 34.",
		},
		{
			name: "band3 selisih geometri",
			p:    params{kind: kindDiffGeom, n: 5, a: 2, d: 1, r: 3},
			stem: "Perhatikan deret angka berikut: 2, 3, 6, 15, 42, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "123",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"69"},
				MisRasioJadiBeda:    {"72"},
				MisSalahHitungDekat: {"122", "124"},
			},
			stepPart: "Suku berikutnya = 42 + 81 = 123.",
		},
		{
			name: "band3 selisih tingkat dua",
			p:    params{kind: kindOrder3, n: 6, a: 1, d: 1, c: 0, e: 1},
			stem: "Perhatikan deret angka berikut: 1, 2, 3, 5, 9, 16, ... Angka berikutnya dalam deret itu adalah ...",
			key:  "27",
			wantMis: map[string][]string{
				MisBedaTerakhir:     {"23"},
				MisSalahHitungDekat: {"26", "28", "29"},
			},
			stepPart: "Suku berikutnya = 16 + 11 = 27.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := build(c.p, 7)
			if d.Stem != c.stem {
				t.Errorf("stem = %q", d.Stem)
			}
			key, got := byMisconception(d)
			if key != c.key {
				t.Errorf("key = %q, want %q", key, c.key)
			}
			if !reflect.DeepEqual(got, c.wantMis) {
				t.Errorf("distractors = %v, want %v", got, c.wantMis)
			}
			if !strings.Contains(strings.Join(d.Steps, "\n"), c.stepPart) {
				t.Errorf("steps %q lack %q", d.Steps, c.stepPart)
			}
		})
	}
}

func TestMisconceptionLabelsCoverEmittedIds(t *testing.T) {
	fam := Family{}
	labels := fam.Misconceptions()
	emitted := map[string]bool{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 500; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			for _, o := range d.Options {
				if o.MisconceptionID == "" {
					continue
				}
				emitted[o.MisconceptionID] = true
				if labels[o.MisconceptionID] == "" {
					t.Fatalf("band %d seed %d: emitted id %q has no label", band, seed, o.MisconceptionID)
				}
			}
		}
	}
	for id := range labels {
		if !emitted[id] {
			t.Errorf("label %q is never emitted", id)
		}
	}
	if fam.ID() != PatternID || fam.Label() != "Deret angka" {
		t.Errorf("id/label = %q/%q", fam.ID(), fam.Label())
	}
}

func TestGenerateErrorsWithoutValidParams(t *testing.T) {
	if _, err := generate(1, gen.Band2, maxAttempts, func(params) bool { return false }); err == nil {
		t.Fatal("expected error when no params are accepted")
	}
	if _, err := (Family{}).Generate(1, gen.Band(9)); err == nil {
		t.Fatal("expected error for unknown band")
	}
}

func TestValidRejectsOutOfRangeTerms(t *testing.T) {
	bad := []params{
		{kind: kindArith, n: 5, a: 3, d: -2},           // goes below 1
		{kind: kindGeom, n: 6, a: 6, r: 3},             // exceeds maxTerm
		{kind: kindAlt, n: 6, a: 5, d: 3, b: 5, d2: 3}, // identical rules
		{kind: kindStepped, n: 5, a: 1, d: 2, e: 0},    // not stepped
		{kind: kindCube, n: 6, s: 7, c: 0},             // 13^3 > maxTerm
	}
	for _, p := range bad {
		if valid(p) {
			t.Errorf("valid(%+v) = true", p)
		}
	}
}
