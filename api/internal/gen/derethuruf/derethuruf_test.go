package derethuruf

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/gentest"
)

var seriesRe = regexp.MustCompile(`^Lanjutkan deret berikut: (.+), \.\.\. Suku berikutnya adalah \.\.\.$`)

type tok struct {
	letter bool
	v      int
}

func parse(t *testing.T, s string) tok {
	t.Helper()
	if n, err := strconv.Atoi(s); err == nil {
		return tok{false, n}
	}
	if len(s) != 1 || s[0] < 'A' || s[0] > 'Z' {
		t.Fatalf("bad token %q", s)
	}
	return tok{true, int(s[0]-'A') + 1}
}

// newton extrapolates one subsequence with Newton forward differences, a
// different route from the generator's difference-table walk: it fits a
// polynomial of the given degree through the last degree+1 points and
// evaluates it one step ahead via binomial weights.
func newton(ys []int, degree int) (int, bool) {
	if len(ys) < degree+2 {
		return 0, false
	}
	// every degree+1 consecutive window must agree on the same prediction.
	pred := func(w []int) int {
		// Lagrange form of the unique degree-d polynomial through w at x=d+1.
		d := len(w) - 1
		total := 0
		for k := 0; k <= d; k++ {
			c := 1
			for j := 0; j <= d; j++ {
				if j != k {
					c *= d + 1 - j
				}
			}
			den := 1
			for j := 0; j <= d; j++ {
				if j != k {
					den *= k - j
				}
			}
			total += w[k] * c / den
		}
		return total
	}
	for i := 1; i+degree < len(ys); i++ {
		// the polynomial through window i must reproduce the next listed value.
		if pred(ys[i-1:i+degree]) != ys[i+degree] {
			return 0, false
		}
	}
	return pred(ys[len(ys)-degree-1:]), true
}

func expected(t *testing.T, stem string) string {
	t.Helper()
	m := seriesRe.FindStringSubmatch(stem)
	if m == nil {
		t.Fatalf("unrecognised stem %q", stem)
	}
	var ts []tok
	for _, s := range strings.Split(m[1], ", ") {
		ts = append(ts, parse(t, s))
	}
	var answers []tok
	for _, h := range [][2]int{{1, 1}, {1, 2}, {2, 1}} {
		period, degree := h[0], h[1]
		var ys []int
		kindOK := true
		for i := len(ts) % period; i < len(ts); i += period {
			ys = append(ys, ts[i].v)
			if ts[i].letter != ts[len(ts)%period].letter {
				kindOK = false
			}
		}
		if !kindOK {
			continue
		}
		if v, ok := newton(ys, degree); ok {
			answers = append(answers, tok{ts[len(ts)%period].letter, v})
		}
	}
	if len(answers) == 0 {
		t.Fatalf("no hypothesis fits %q", stem)
	}
	for _, a := range answers {
		if a != answers[0] {
			t.Fatalf("ambiguous series %q: %v", stem, answers)
		}
	}
	a := answers[0]
	if a.letter {
		return string(rune('A' + a.v - 1))
	}
	return strconv.Itoa(a.v)
}

func TestKeyMatchesIndependentExtrapolation(t *testing.T) {
	gentest.Run(t, Family{}, 5000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		if got, want := gentest.KeyText(d), expected(t, d.Stem); got != want {
			t.Fatalf("key %q, independent %q", got, want)
		}
		for _, o := range d.Options {
			tk := parse(t, o.Text)
			if tk.letter && (tk.v < 1 || tk.v > 26) || tk.v < 1 {
				t.Fatalf("option %q out of range", o.Text)
			}
		}
	})
}

func TestBandsEmitExpectedShapes(t *testing.T) {
	hasNumber := regexp.MustCompile(`: .*\b\d+\b.*, \.\.\.`)
	numbers := 0
	for s := int64(0); s < 300; s++ {
		d, _ := Family{}.Generate(s, gen.Band2)
		if hasNumber.MatchString(d.Stem) {
			numbers++
		}
		if d1, _ := (Family{}).Generate(s, gen.Band1); hasNumber.MatchString(d1.Stem) {
			t.Fatalf("band 1 must be letters only: %q", d1.Stem)
		}
	}
	if numbers == 0 || numbers == 300 {
		t.Fatalf("band 2 should mix letter-only and letter-number series, got %d/300", numbers)
	}
}

func TestBadBand(t *testing.T) {
	if _, err := (Family{}).Generate(1, 7); err == nil {
		t.Fatal("want error for unknown band")
	}
}
