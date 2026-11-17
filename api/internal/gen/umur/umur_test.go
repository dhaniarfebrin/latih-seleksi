package umur

import (
	"regexp"
	"strconv"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/gentest"
)

var (
	sumDiffRe = regexp.MustCompile(`^Jumlah umur (\w+) dan (\w+) adalah (\d+) tahun, sedangkan selisih umur mereka (\d+) tahun\. Umur (\w+) sekarang`)
	pastRe    = regexp.MustCompile(`^Umur (\w+) sekarang (\d+) kali umur (\w+)\. (\d+) tahun yang lalu, jumlah umur mereka adalah (\d+) tahun\. Umur (\w+) sekarang`)
	futureRe  = regexp.MustCompile(`^(\d+) tahun yang lalu, umur (\w+) adalah (\d+) kali umur (\w+)\. (\d+) tahun lagi, jumlah umur mereka menjadi (\d+) tahun\. Umur (\w+) sekarang`)
)

func n(s string) int { v, _ := strconv.Atoi(s); return v }

// solve brute-forces every (older, younger) age pair against the constraints
// read from the stem and returns the asked age; it fails unless exactly one
// pair fits.
func solve(t *testing.T, stem string) int {
	t.Helper()
	var fits func(a, b int) bool
	var asked, older string
	if m := sumDiffRe.FindStringSubmatch(stem); m != nil {
		s, d := n(m[3]), n(m[4])
		older, asked = m[1], m[5]
		fits = func(a, b int) bool { return a+b == s && a-b == d }
	} else if m := pastRe.FindStringSubmatch(stem); m != nil {
		k, y, s := n(m[2]), n(m[4]), n(m[5])
		older, asked = m[1], m[6]
		fits = func(a, b int) bool { return a == k*b && (a-y)+(b-y) == s && b-y > 0 }
	} else if m := futureRe.FindStringSubmatch(stem); m != nil {
		x, mult, y, s := n(m[1]), n(m[3]), n(m[5]), n(m[6])
		older, asked = m[2], m[7]
		fits = func(a, b int) bool { return a-x == mult*(b-x) && b-x > 0 && (a+y)+(b+y) == s }
	} else {
		t.Fatalf("unrecognised stem %q", stem)
	}
	var sols []int
	for a := 1; a <= 200; a++ {
		for b := 1; b <= 200; b++ {
			if fits(a, b) {
				if asked == older {
					sols = append(sols, a)
				} else {
					sols = append(sols, b)
				}
			}
		}
	}
	if len(sols) != 1 {
		t.Fatalf("want exactly one solution, got %v for %q", sols, stem)
	}
	return sols[0]
}

func TestKeyMatchesBruteForce(t *testing.T) {
	re := regexp.MustCompile(`^(\d+) tahun$`)
	gentest.Run(t, Family{}, 5000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		m := re.FindStringSubmatch(gentest.KeyText(d))
		if m == nil {
			t.Fatalf("key text %q", gentest.KeyText(d))
		}
		if got, want := n(m[1]), solve(t, d.Stem); got != want {
			t.Fatalf("key %d, brute force %d", got, want)
		}
		wantVariant := map[gen.Band]*regexp.Regexp{gen.Band1: sumDiffRe, gen.Band2: pastRe, gen.Band3: futureRe}[band]
		if !wantVariant.MatchString(d.Stem) {
			t.Fatalf("band %d emitted the wrong variant", band)
		}
	})
}

func TestBadBand(t *testing.T) {
	if _, err := (Family{}).Generate(1, 9); err == nil {
		t.Fatal("want error for unknown band")
	}
}

// TestAgeGapIsBelievable recovers both ages from the stem and the key and
// checks the older relative is at least the pair's gap older.
func TestAgeGapIsBelievable(t *testing.T) {
	gap := map[string]int{}
	for _, p := range pairs {
		gap[p.older] = p.gap
	}
	gentest.Run(t, Family{}, 2000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		key := n(regexp.MustCompile(`^(\d+)`).FindString(gentest.KeyText(d)))
		var older string
		var diff int
		switch band {
		case gen.Band1:
			m := sumDiffRe.FindStringSubmatch(d.Stem)
			older, diff = m[1], n(m[4])
		case gen.Band2:
			m := pastRe.FindStringSubmatch(d.Stem)
			older = m[1]
			k := n(m[2])
			if m[6] == older {
				diff = key - key/k
			} else {
				diff = key*k - key
			}
		default:
			m := futureRe.FindStringSubmatch(d.Stem)
			older = m[2]
			other := n(m[6]) - 2*n(m[5]) - key
			diff = key - other
			if m[7] != older {
				diff = -diff
			}
		}
		if diff < gap[older] {
			t.Fatalf("%s is only %d years older, want >= %d: %q", older, diff, gap[older], d.Stem)
		}
	})
}
