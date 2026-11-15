package insight

import (
	"context"
	"reflect"
	"testing"
	"time"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func fact(attempt, pattern string, correct bool, ms int) Fact {
	return Fact{AttemptID: attempt, Mode: "drill", StartedAt: t0, Section: "TIU", PatternID: pattern, Correct: correct, TimeMS: ms}
}

func TestEmptyIsNonNil(t *testing.T) {
	d := Build(nil, nil)
	if d.Sections == nil || d.Patterns == nil || d.Sessions == nil || d.Queue == nil {
		t.Fatalf("nil slice in %+v", d)
	}
}

func TestAccuracyAndAverageTime(t *testing.T) {
	facts := []Fact{fact("a", "p", true, 1000), fact("a", "p", false, 3000), fact("a", "p", true, 2000)}
	d := Build(facts, nil)
	p := d.Patterns[0]
	if p.Answered != 3 || p.Correct != 2 || p.AvgTimeMS != 2000 {
		t.Errorf("pattern = %+v", p)
	}
	s := d.Sections[0]
	if s.Section != "TIU" || s.Answered != 3 || s.Correct != 2 || s.AvgTimeMS != 2000 {
		t.Errorf("section = %+v", s)
	}
}

func TestMedian(t *testing.T) {
	if got := median([]int{5, 1, 3}); got != 3 {
		t.Errorf("odd median = %v", got)
	}
	if got := median([]int{4, 1, 2, 3}); got != 2.5 {
		t.Errorf("even median = %v", got)
	}
	if got := median(nil); got != 0 {
		t.Errorf("empty median = %v", got)
	}
}

func TestCauseLabels(t *testing.T) {
	correct := []Fact{fact("a", "p", true, 1000), fact("a", "p", true, 1000), fact("a", "p", true, 1000)} // median 1000
	cases := []struct {
		name string
		ms   int
		want CauseCounts
	}{
		{"slow", 1501, CauseCounts{TimePressure: 1}},
		{"fast", 499, CauseCounts{Careless: 1}},
		{"middle", 1000, CauseCounts{Misunderstood: 1}},
		{"exactly 1.5x is not slow", 1500, CauseCounts{Misunderstood: 1}},
		{"exactly 0.5x is not fast", 500, CauseCounts{Misunderstood: 1}},
	}
	for _, c := range cases {
		d := Build(append(append([]Fact{}, correct...), fact("a", "p", false, c.ms)), nil)
		if got := d.Patterns[0].Causes; got != c.want {
			t.Errorf("%s: causes = %+v, want %+v", c.name, got, c.want)
		}
	}
}

func TestFewCorrectIsUnclear(t *testing.T) {
	facts := []Fact{fact("a", "p", true, 1000), fact("a", "p", true, 1000), fact("a", "p", false, 9000)}
	if got := Build(facts, nil).Patterns[0].Causes; got != (CauseCounts{Unclear: 1}) {
		t.Errorf("causes = %+v", got)
	}
}

func TestTopMisconception(t *testing.T) {
	w := func(m string) Fact { f := fact("a", "p", false, 1); f.MisconceptionID = m; return f }
	got := Build([]Fact{w("b"), w("a"), w("b"), w("a"), w("")}, nil).Patterns[0].TopMisconceptionID
	if got != "a" {
		t.Errorf("top = %q, want a (tie to smallest id)", got)
	}
	if got := Build([]Fact{fact("a", "p", true, 1)}, nil).Patterns[0].TopMisconceptionID; got != "" {
		t.Errorf("top = %q, want empty", got)
	}
}

func repeat(pattern string, correct, wrong int) []Fact {
	var out []Fact
	for i := 0; i < correct; i++ {
		out = append(out, fact("a", pattern, true, 1000))
	}
	for i := 0; i < wrong; i++ {
		out = append(out, fact("a", pattern, false, 1000))
	}
	return out
}

func TestQueueOrderAndUntriedFill(t *testing.T) {
	var facts []Fact
	facts = append(facts, repeat("mid", 3, 1)...)  // 75%
	facts = append(facts, repeat("low", 1, 3)...)  // 25%
	facts = append(facts, repeat("good", 4, 1)...) // 80%: not queued
	facts = append(facts, repeat("few", 0, 2)...)  // 0% but only 2 answers
	d := Build(facts, []string{"good", "new1", "low", "new2"})
	want := []QueueItem{
		{PatternID: "low", Reason: ReasonLowAccuracy, AccuracyPct: 25},
		{PatternID: "mid", Reason: ReasonLowAccuracy, AccuracyPct: 75},
		{PatternID: "new1", Reason: ReasonUntried},
	}
	if !reflect.DeepEqual(d.Queue, want) {
		t.Errorf("queue = %+v, want %+v", d.Queue, want)
	}
}

// The queue percent is rounded like the UI's own percent (1/8 = 12.5% shows as 13%), not truncated.
func TestQueueAccuracyRounds(t *testing.T) {
	d := Build(repeat("p", 1, 7), nil)
	if len(d.Queue) != 1 || d.Queue[0].AccuracyPct != 13 {
		t.Errorf("queue = %+v, want one item at 13%%", d.Queue)
	}
}

func TestQueueCappedAtThree(t *testing.T) {
	var facts []Fact
	for _, p := range []string{"a", "b", "c", "d"} {
		facts = append(facts, repeat(p, 0, 3)...)
	}
	if n := len(Build(facts, []string{"x"}).Queue); n != MaxQueue {
		t.Errorf("queue length = %d", n)
	}
}

func TestSessionsNewestFirstAndCapped(t *testing.T) {
	var facts []Fact
	for i := 0; i < 12; i++ {
		f := fact(string(rune('a'+i)), "p", i%2 == 0, 1)
		f.StartedAt = t0.Add(time.Duration(i) * time.Hour)
		facts = append(facts, f)
	}
	d := Build(facts, nil)
	if len(d.Sessions) != MaxSessions {
		t.Fatalf("sessions = %d", len(d.Sessions))
	}
	if d.Sessions[0].AttemptID != "l" || d.Sessions[9].AttemptID != "c" {
		t.Errorf("order = %s ... %s", d.Sessions[0].AttemptID, d.Sessions[9].AttemptID)
	}
}

// A drill attempt has one pattern and names it; a simulation mixes patterns and names none.
func TestSessionPatternOnlyForDrill(t *testing.T) {
	sim := fact("sim", "p1", true, 1)
	sim.Mode = "simulasi"
	sim.StartedAt = t0.Add(time.Hour)
	sim2 := fact("sim", "p2", false, 1)
	sim2.Mode = "simulasi"
	sim2.StartedAt = t0.Add(time.Hour)
	d := Build([]Fact{fact("dr", "p1", true, 1), sim, sim2}, nil)
	if len(d.Sessions) != 2 {
		t.Fatalf("sessions = %+v", d.Sessions)
	}
	if d.Sessions[0].AttemptID != "sim" || d.Sessions[0].PatternID != "" {
		t.Errorf("simulation session = %+v, want no pattern", d.Sessions[0])
	}
	if d.Sessions[1].AttemptID != "dr" || d.Sessions[1].PatternID != "p1" {
		t.Errorf("drill session = %+v, want pattern p1", d.Sessions[1])
	}
}

type fakeStore struct {
	facts []Fact
	limit int
}

func (f *fakeStore) ResponseFacts(_ context.Context, _ string, limit int) ([]Fact, error) {
	f.limit = limit
	return f.facts, nil
}

func TestServiceUsesFactLimit(t *testing.T) {
	st := &fakeStore{facts: repeat("p", 1, 0)}
	d, err := New(st, nil).Dashboard(context.Background(), "u")
	if err != nil || st.limit != FactLimit || len(d.Patterns) != 1 {
		t.Errorf("limit=%d err=%v d=%+v", st.limit, err, d)
	}
}
