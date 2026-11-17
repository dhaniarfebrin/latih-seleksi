package himpunan

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/gentest"
)

var (
	intRe  = regexp.MustCompile(`\d+`)
	onlyRe = regexp.MustCompile(`, \d+ gemar (.+?), \d+ gemar (.+?), dan \d+ gemar keduanya\. Banyak \w+ yang hanya gemar (.+?) adalah`)
)

func ints(s string) []int {
	var out []int
	for _, m := range intRe.FindAllString(s, -1) {
		v, _ := strconv.Atoi(m)
		out = append(out, v)
	}
	return out
}

// expected applies inclusion-exclusion to the numbers read from the stem.
func expected(t *testing.T, stem string, band gen.Band) int {
	t.Helper()
	v := ints(stem)
	switch band {
	case gen.Band1:
		n, a, b, c := v[0], v[1], v[2], v[3]
		union := a + b - c
		switch {
		case strings.Contains(stem, "tidak gemar satu pun"):
			return n - union
		case strings.Contains(stem, "paling sedikit satu"):
			return union
		}
		m := onlyRe.FindStringSubmatch(stem)
		if m == nil {
			t.Fatalf("unrecognised question %q", stem)
		}
		if m[3] == m[1] {
			return a - c
		}
		return b - c
	case gen.Band2:
		n, a, b, none := v[0], v[1], v[2], v[3]
		return a + b - (n - none)
	}
	n, a, b, c, ab, ac, bc, abc := v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7]
	union := a + b + c - ab - ac - bc + abc
	switch {
	case strings.Contains(stem, "tidak gemar satu pun"):
		return n - union
	case strings.Contains(stem, "paling sedikit satu"):
		return union
	case strings.Contains(stem, "tepat satu"):
		return a + b + c - 2*(ab+ac+bc) + 3*abc
	}
	t.Fatalf("unrecognised question %q", stem)
	return 0
}

func TestKeyMatchesInclusionExclusion(t *testing.T) {
	gentest.Run(t, Family{}, 5000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		got, err := strconv.Atoi(gentest.KeyText(d))
		if err != nil || got < 0 {
			t.Fatalf("key text %q", gentest.KeyText(d))
		}
		if want := expected(t, d.Stem, band); got != want {
			t.Fatalf("key %d, independent %d", got, want)
		}
	})
}

func TestBadBand(t *testing.T) {
	if _, err := (Family{}).Generate(1, 4); err == nil {
		t.Fatal("want error for unknown band")
	}
}
