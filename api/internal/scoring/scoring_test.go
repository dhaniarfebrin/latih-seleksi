package scoring

import (
	"fmt"
	"strings"
	"testing"

	"latihseleksi/api/internal/spec"
)

func load(t *testing.T, id string) spec.Spec {
	t.Helper()
	r, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	s, ok := r.Get(id)
	if !ok {
		t.Fatalf("spec %s missing", id)
	}
	return s
}

// attempt builds responses: k correct out of n for TWK and TIU (the rest
// wrong), and the given TKP weights (nil = all unanswered).
type attempt struct {
	twk, tiu int
	tkp      []int // 45 weights; 0 = unanswered
}

func (a attempt) responses() []Response {
	var rs []Response
	for i := 0; i < 30; i++ {
		rs = append(rs, Response{Section: "TWK", Answered: true, Correct: i < a.twk})
	}
	for i := 0; i < 35; i++ {
		rs = append(rs, Response{Section: "TIU", Answered: true, Correct: i < a.tiu})
	}
	for i := 0; i < 45; i++ {
		w := 0
		if a.tkp != nil {
			w = a.tkp[i]
		}
		rs = append(rs, Response{Section: "TKP", Answered: w != 0, Weight: w})
	}
	return rs
}

// tkpSum returns 45 weights summing to total (each 1..5, total in 45..225).
func tkpSum(total int) []int {
	w := make([]int, 45)
	for i := range w {
		w[i] = 1
	}
	rest := total - 45
	for i := 0; rest > 0; i = (i + 1) % 45 {
		if w[i] < 5 {
			w[i]++
			rest--
		}
	}
	return w
}

func sectionScore(t *testing.T, r Result, id string) SectionResult {
	t.Helper()
	for _, s := range r.Sections {
		if s.ID == id {
			return s
		}
	}
	t.Fatalf("no section %s in result", id)
	return SectionResult{}
}

func TestRightWrongSections(t *testing.T) {
	s := load(t, "skd-2024")
	// Oracle: hand-written, 5 points per correct item, 0 otherwise.
	for k := 0; k <= 30; k++ {
		res, err := Score(s, "umum", attempt{twk: k}.responses())
		if err != nil {
			t.Fatal(err)
		}
		if got := sectionScore(t, res, "TWK").Score; got != 5*k {
			t.Errorf("TWK k=%d: score %d, want %d", k, got, 5*k)
		}
	}
	for k := 0; k <= 35; k++ {
		res, err := Score(s, "umum", attempt{tiu: k}.responses())
		if err != nil {
			t.Fatal(err)
		}
		if got := sectionScore(t, res, "TIU").Score; got != 5*k {
			t.Errorf("TIU k=%d: score %d, want %d", k, got, 5*k)
		}
	}
}

func TestBlankRightWrongScoresZero(t *testing.T) {
	s := load(t, "skd-2024")
	rs := attempt{twk: 30, tiu: 35}.responses()
	for i := range rs {
		if rs[i].Section == "TWK" && i < 10 {
			rs[i] = Response{Section: "TWK"} // unanswered, even though "Correct" is false
		}
	}
	res, err := Score(s, "umum", rs)
	if err != nil {
		t.Fatal(err)
	}
	if got := sectionScore(t, res, "TWK").Score; got != 100 {
		t.Errorf("TWK with 10 blank of 30: %d, want 100", got)
	}
}

func TestTKPWeights(t *testing.T) {
	s := load(t, "skd-2024")
	uniform := func(w int) []int {
		out := make([]int, 45)
		for i := range out {
			out[i] = w
		}
		return out
	}
	alt := make([]int, 45) // 1,2,3,4,5 repeating: nine rounds of 15
	for i := range alt {
		alt[i] = i%5 + 1
	}
	cases := []struct {
		name string
		tkp  []int
		want int
	}{
		{"all unanswered", nil, 0},
		{"all 1", uniform(1), 45},
		{"all 2", uniform(2), 90},
		{"all 3", uniform(3), 135},
		{"all 4", uniform(4), 180},
		{"all 5", uniform(5), 225},
		{"cycle 1..5", alt, 9 * 15},
		{"sum 166", tkpSum(166), 166},
	}
	for _, c := range cases {
		res, err := Score(s, "umum", attempt{tkp: c.tkp}.responses())
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got := sectionScore(t, res, "TKP").Score; got != c.want {
			t.Errorf("%s: TKP %d, want %d", c.name, got, c.want)
		}
		if got := sectionScore(t, res, "TKP").Max; got != 225 {
			t.Errorf("%s: TKP max %d, want 225", c.name, got)
		}
	}
	// Some answered, some blank: blank contributes 0.
	mixed := uniform(5)
	for i := 0; i < 15; i++ {
		mixed[i] = 0
	}
	res, err := Score(s, "umum", attempt{tkp: mixed}.responses())
	if err != nil {
		t.Fatal(err)
	}
	if got := sectionScore(t, res, "TKP").Score; got != 150 {
		t.Errorf("30 x 5 + 15 blank: %d, want 150", got)
	}
}

// Boundaries are written by hand from the decrees: score n-1 fails, n passes.
func TestThresholdBoundaries(t *testing.T) {
	// perfect attempt except for the one section under test.
	full := attempt{twk: 30, tiu: 35, tkp: tkpSum(225)}
	cases := []struct {
		spec, track string
		section     string
		n           int
	}{
		{"skd-2024", "umum", "TWK", 65},
		{"skd-2024", "umum", "TIU", 80},
		{"skd-2024", "umum", "TKP", 166},
		{"skd-sekolah-kedinasan-2026", "umum", "TWK", 65},
		{"skd-sekolah-kedinasan-2026", "umum", "TIU", 80},
		{"skd-sekolah-kedinasan-2026", "umum", "TKP", 156},
		{"skd-2024", "cumlaude", "TIU", 85},
		{"skd-2024", "afirmasi", "TIU", 60},
		{"skd-sekolah-kedinasan-2026", "afirmasi", "TIU", 55},
	}
	for _, c := range cases {
		s := load(t, c.spec)
		// TWK and TIU move in steps of 5, so n-1 is reached as n-5 (score/5 rounds down).
		for _, score := range []int{c.n - 1, c.n} {
			a := full
			switch c.section {
			case "TWK":
				a.twk = score / 5
			case "TIU":
				a.tiu = score / 5
			case "TKP":
				a.tkp = tkpSum(score)
			}
			res, err := Score(s, c.track, a.responses())
			if err != nil {
				t.Fatalf("%s/%s/%s: %v", c.spec, c.track, c.section, err)
			}
			sec := sectionScore(t, res, c.section)
			wantMet := sec.Score >= c.n
			if sec.Met != wantMet || sec.Min == nil || *sec.Min != c.n {
				t.Errorf("%s/%s %s=%d: met=%v min=%v, want met=%v min=%d", c.spec, c.track, c.section, sec.Score, sec.Met, sec.Min, wantMet, c.n)
			}
			// the totals are high in all these cases, so Passed follows the section alone.
			if res.Passed != wantMet {
				t.Errorf("%s/%s %s=%d: passed=%v, want %v", c.spec, c.track, c.section, sec.Score, res.Passed, wantMet)
			}
		}
	}
}

func TestCumulativeBoundaries(t *testing.T) {
	cases := []struct {
		spec, track string
		n           int
		twk, tiu    int // correct items; TIU is enough to clear its own threshold
	}{
		{"skd-2024", "cumlaude", 311, 4, 17},                   // TIU 85
		{"skd-2024", "afirmasi", 286, 2, 12},                   // TIU 60
		{"skd-sekolah-kedinasan-2026", "afirmasi", 281, 2, 11}, // TIU 55
	}
	for _, c := range cases {
		s := load(t, c.spec)
		for _, total := range []int{c.n - 1, c.n} {
			// TWK and TIU are fixed, TKP carries the rest of the total.
			tkp := total - 5*c.twk - 5*c.tiu
			if tkp < 45 || tkp > 225 {
				t.Fatalf("test setup: TKP %d out of range", tkp)
			}
			res, err := Score(s, c.track, attempt{twk: c.twk, tiu: c.tiu, tkp: tkpSum(tkp)}.responses())
			if err != nil {
				t.Fatal(err)
			}
			want := total >= c.n
			if res.Total != total || res.TotalMet != want || res.Passed != want {
				t.Errorf("%s/%s total=%d: got total=%d met=%v passed=%v, want met=%v", c.spec, c.track, total, res.Total, res.TotalMet, res.Passed, want)
			}
			if res.TotalMin == nil || *res.TotalMin != c.n {
				t.Errorf("%s/%s: TotalMin = %v, want %d", c.spec, c.track, res.TotalMin, c.n)
			}
		}
	}
}

// A special track can clear the cumulative total yet fail TIU.
func TestSpecialTrackFailsTIU(t *testing.T) {
	s := load(t, "skd-2024")
	// TWK 150 + TIU 50 (10 correct) + TKP 225 = 425 >= 311, but TIU 50 < 85.
	res, err := Score(s, "cumlaude", attempt{twk: 30, tiu: 10, tkp: tkpSum(225)}.responses())
	if err != nil {
		t.Fatal(err)
	}
	if !res.TotalMet || sectionScore(t, res, "TIU").Met || res.Passed {
		t.Errorf("total met=%v TIU met=%v passed=%v, want true/false/false", res.TotalMet, sectionScore(t, res, "TIU").Met, res.Passed)
	}
}

// Sections without a threshold in the track cannot fail it.
func TestUnthresholdedSectionCannotFail(t *testing.T) {
	s := load(t, "skd-2024")
	// afirmasi: only TIU >= 60 and total >= 286. TWK 0 and TKP 45 still pass.
	res, err := Score(s, "afirmasi", attempt{tiu: 35, tkp: tkpSum(225)}.responses())
	if err != nil {
		t.Fatal(err)
	}
	twk := sectionScore(t, res, "TWK")
	if twk.Min != nil || !twk.Met || !res.Passed {
		t.Errorf("TWK min=%v met=%v passed=%v", twk.Min, twk.Met, res.Passed)
	}
	if res.SpecID != "skd-2024" || res.Track != "afirmasi" || res.MaxTotal != 550 || res.Total != 175+225 {
		t.Errorf("result header %+v", res)
	}
}

func TestScoreErrors(t *testing.T) {
	s := load(t, "skd-2024")
	good := attempt{twk: 1}.responses()
	cases := []struct {
		name  string
		track string
		rs    []Response
		want  string
	}{
		{"unknown track", "nope", good, "unknown track"},
		{"unknown section", "umum", append(append([]Response(nil), good...), Response{Section: "XYZ"}), "unknown section"},
		{"too few TWK", "umum", good[1:], "section TWK: got 29"},
		{"too many TIU", "umum", append(append([]Response(nil), good...), Response{Section: "TIU"}), "section TIU: got 36"},
		{"no responses", "umum", nil, "got 0"},
		{"weight too high", "umum", withWeight(good, 6), "weight 6"},
		{"weight negative", "umum", withWeight(good, -1), "weight -1"},
	}
	for _, c := range cases {
		_, err := Score(s, c.track, c.rs)
		if err == nil {
			t.Errorf("%s: no error", c.name)
			continue
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %q does not mention %q", c.name, err, c.want)
		}
	}
	// weight 0 with Answered=true is not a valid weight either.
	if _, err := Score(s, "umum", answeredWeight0(good)); err == nil {
		t.Error("answered TKP item with weight 0 accepted")
	}
}

func withWeight(rs []Response, w int) []Response {
	out := append([]Response(nil), rs...)
	for i := range out {
		if out[i].Section == "TKP" {
			out[i] = Response{Section: "TKP", Answered: true, Weight: w}
			break
		}
	}
	return out
}

func answeredWeight0(rs []Response) []Response {
	return withWeight(rs, 0)
}

func TestSectionOrderFollowsSpec(t *testing.T) {
	s := load(t, "skd-2024")
	rs := attempt{}.responses()
	// shuffle: reverse the responses; order of input must not matter.
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	res, err := Score(s, "umum", rs)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, sec := range res.Sections {
		ids = append(ids, sec.ID)
	}
	if got := fmt.Sprint(ids); got != "[TWK TIU TKP]" {
		t.Errorf("section order %s", got)
	}
}
