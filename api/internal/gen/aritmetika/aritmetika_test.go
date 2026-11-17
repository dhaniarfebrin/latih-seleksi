package aritmetika

import (
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/gentest"
)

var (
	moneyRe  = regexp.MustCompile(`Rp ([\d.]+)`)
	pctRe    = regexp.MustCompile(`(\d+)%`)
	pctOptRe = regexp.MustCompile(`^(\d+)%$`)
)

func amounts(s string) []int64 {
	var out []int64
	for _, m := range moneyRe.FindAllStringSubmatch(s, -1) {
		v, _ := strconv.ParseInt(strings.ReplaceAll(m[1], ".", ""), 10, 64)
		out = append(out, v)
	}
	return out
}

func pcts(s string) []int64 {
	var out []int64
	for _, m := range pctRe.FindAllStringSubmatch(s, -1) {
		v, _ := strconv.ParseInt(m[1], 10, 64)
		out = append(out, v)
	}
	return out
}

func intsAfter(s, re string) int64 {
	m := regexp.MustCompile(re).FindStringSubmatch(s)
	if m == nil {
		return -1
	}
	v, _ := strconv.ParseInt(m[1], 10, 64)
	return v
}

// expected recomputes the key from the stem with exact rationals.
func expected(t *testing.T, stem string) string {
	t.Helper()
	a, p := amounts(stem), pcts(stem)
	r := func(v int64) *big.Rat { return big.NewRat(v, 1) }
	frac := func(base *big.Rat, pc int64) *big.Rat { return new(big.Rat).Mul(base, big.NewRat(pc, 100)) }
	var res *big.Rat
	money := true
	switch {
	case strings.Contains(stem, "dengan laba"):
		res = frac(r(a[0]), 100+p[0])
	case strings.Contains(stem, "dengan rugi"):
		res = frac(r(a[0]), 100-p[0])
	case strings.Contains(stem, "sebelum diskon adalah") && strings.Contains(stem, "promo"):
		res = frac(r(a[0]), 100-p[0])
	case strings.Contains(stem, "bunga tunggal"):
		y := intsAfter(stem, `Setelah (\d+) tahun`)
		res = new(big.Rat).Add(r(a[0]), frac(r(a[0]), p[0]*y))
	case strings.Contains(stem, "Persentase"):
		money = false
		diff := new(big.Rat).Sub(r(a[1]), r(a[0]))
		diff.Abs(diff)
		res = frac(new(big.Rat).Quo(diff, r(a[0])), 100*100)
	case strings.Contains(stem, "diskon tambahan"):
		res = frac(frac(r(a[0]), 100-p[0]), 100-p[1])
	case strings.Contains(stem, "dikenai pajak"):
		res = frac(frac(r(a[0]), 100-p[0]), 100+p[1])
	case strings.Contains(stem, "sebelum diskon adalah"):
		res = new(big.Rat).Quo(r(a[0]), big.NewRat(100-p[0], 100))
	default:
		t.Fatalf("unrecognised stem %q", stem)
	}
	if !res.IsInt() {
		t.Fatalf("non-integer key %s for %q", res.RatString(), stem)
	}
	if !money {
		return res.Num().String() + "%"
	}
	return "Rp " + dots(res.Num().String())
}

func dots(s string) string {
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	return strings.Join(append([]string{s}, parts...), ".")
}

func TestKeyMatchesIndependentCalculation(t *testing.T) {
	gentest.Run(t, Family{}, 5000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		if got, want := gentest.KeyText(d), expected(t, d.Stem); got != want {
			t.Fatalf("key %q, independent %q", got, want)
		}
		for _, o := range d.Options {
			if m := pctOptRe.FindStringSubmatch(o.Text); m != nil && intsAfter(o.Text, `^(\d+)%$`) <= 0 {
				t.Fatalf("non-positive percent option %q", o.Text)
			}
		}
	})
}

func TestBandsEmitExpectedVariants(t *testing.T) {
	seen := map[gen.Band]map[string]bool{gen.Band1: {}, gen.Band2: {}, gen.Band3: {}}
	for s := int64(0); s < 400; s++ {
		for b := range seen {
			d, err := Family{}.Generate(s, b)
			if err != nil {
				t.Fatal(err)
			}
			for _, k := range []string{"dengan laba", "dengan rugi", "promo", "bunga tunggal", "Persentase", "diskon tambahan", "dikenai pajak", "menjadi Rp"} {
				if strings.Contains(d.Stem, k) {
					seen[b][k] = true
				}
			}
		}
	}
	want := map[gen.Band]int{gen.Band1: 3, gen.Band2: 2, gen.Band3: 3}
	for b, w := range want {
		if len(seen[b]) < w {
			t.Errorf("band %d variants seen: %v", b, seen[b])
		}
	}
}

func TestBadBand(t *testing.T) {
	if _, err := (Family{}).Generate(1, 0); err == nil {
		t.Fatal("want error for unknown band")
	}
}
