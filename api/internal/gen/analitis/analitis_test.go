package analitis

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/gentest"
)

var (
	namesRe  = regexp.MustCompile(`^\w+ orang, yaitu (.+?), duduk berjajar`)
	clueRe   = regexp.MustCompile(`\(\d+\) ([^.]+)\.`)
	nameQRe  = regexp.MustCompile(`Siapa yang duduk di urutan ke-(\d+) dari kiri\?$`)
	posQRe   = regexp.MustCompile(`Pada urutan ke berapa dari kiri (\w+) duduk\?$`)
	adjRe    = regexp.MustCompile(`^(\w+) duduk bersebelahan dengan (\w+)$`)
	notAdjRe = regexp.MustCompile(`^(\w+) tidak duduk bersebelahan dengan (\w+)$`)
	leftRe   = regexp.MustCompile(`^(\w+) duduk di sebelah kiri (\w+), boleh berjarak$`)
	immRe    = regexp.MustCompile(`^(\w+) duduk tepat di sebelah kiri (\w+)$`)
	firstRe  = regexp.MustCompile(`^(\w+) duduk paling kiri$`)
	lastRe   = regexp.MustCompile(`^(\w+) duduk paling kanan$`)
)

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// rule is a clue parsed back from the stem text.
type rule func(slot map[string]int, n int) bool

func parseRule(t *testing.T, s string) rule {
	t.Helper()
	if m := adjRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return abs(p[m[1]]-p[m[2]]) == 1 }
	}
	if m := notAdjRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return abs(p[m[1]]-p[m[2]]) > 1 }
	}
	if m := leftRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return p[m[1]] < p[m[2]] }
	}
	if m := immRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return p[m[1]]+1 == p[m[2]] }
	}
	if m := firstRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return p[m[1]] == 0 }
	}
	if m := lastRe.FindStringSubmatch(s); m != nil {
		return func(p map[string]int, n int) bool { return p[m[1]] == n-1 }
	}
	t.Fatalf("unrecognised clue %q", s)
	return nil
}

// solutions enumerates every seating of names with Heap's algorithm and
// returns those that satisfy all rules.
func solutions(names []string, rules []rule) []map[string]int {
	n := len(names)
	arr := append([]string(nil), names...)
	var out []map[string]int
	check := func() {
		p := map[string]int{}
		for i, s := range arr {
			p[s] = i
		}
		for _, r := range rules {
			if !r(p, n) {
				return
			}
		}
		out = append(out, p)
	}
	var heap func(k int)
	heap = func(k int) {
		if k == 1 {
			check()
			return
		}
		heap(k - 1)
		for i := 0; i < k-1; i++ {
			if k%2 == 0 {
				arr[i], arr[k-1] = arr[k-1], arr[i]
			} else {
				arr[0], arr[k-1] = arr[k-1], arr[0]
			}
			heap(k - 1)
		}
	}
	heap(n)
	return out
}

func TestKeyMatchesBruteForce(t *testing.T) {
	clueCounts := map[gen.Band]int{}
	gentest.Run(t, Family{}, 5000, 100, func(t *testing.T, d gen.Draft, band gen.Band) {
		m := namesRe.FindStringSubmatch(d.Stem)
		if m == nil {
			t.Fatalf("no name list in %q", d.Stem)
		}
		names := strings.Split(strings.Replace(m[1], " dan ", ", ", 1), ", ")
		if len(names) != slotsFor(band) {
			t.Fatalf("band %d: %d seats", band, len(names))
		}
		var rules []rule
		for _, c := range clueRe.FindAllStringSubmatch(d.Stem, -1) {
			rules = append(rules, parseRule(t, c[1]))
		}
		clueCounts[band] += len(rules)
		sols := solutions(names, rules)
		if len(sols) != 1 {
			t.Fatalf("want exactly one arrangement, got %d", len(sols))
		}
		sol := sols[0]
		want := ""
		if q := nameQRe.FindStringSubmatch(d.Stem); q != nil {
			k, _ := strconv.Atoi(q[1])
			for name, slot := range sol {
				if slot == k-1 {
					want = name
				}
			}
		} else if q := posQRe.FindStringSubmatch(d.Stem); q != nil {
			want = strconv.Itoa(sol[q[1]] + 1)
		} else {
			t.Fatalf("unrecognised question in %q", d.Stem)
		}
		if got := gentest.KeyText(d); got != want {
			t.Fatalf("key %q, brute force %q", got, want)
		}
	})
	for b, c := range clueCounts {
		if c < 3*5000 {
			t.Errorf("band %d averages fewer than 3 clues", b)
		}
	}
}

func TestBadBand(t *testing.T) {
	if _, err := (Family{}).Generate(1, 8); err == nil {
		t.Fatal("want error for unknown band")
	}
}
