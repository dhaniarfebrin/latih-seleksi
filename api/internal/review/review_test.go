package review

import (
	"context"
	"sort"
	"testing"
	"time"

	"latihseleksi/api/internal/fsrs"
)

type fakeStore struct {
	cards map[string]fsrs.Card
	logs  map[string]Log
}

func newFake() *fakeStore { return &fakeStore{cards: map[string]fsrs.Card{}, logs: map[string]Log{}} }

func (f *fakeStore) Apply(_ context.Context, user, attempt string, outcomes []Outcome, step Step) error {
	for _, o := range outcomes {
		lk := attempt + "/" + string(o.Kind) + "/" + o.Key
		if _, done := f.logs[lk]; done {
			continue
		}
		ck := user + "/" + string(o.Kind) + "/" + o.Key
		var prev *fsrs.Card
		if c, ok := f.cards[ck]; ok {
			prev = &c
		}
		c, l := step(prev, o)
		f.cards[ck], f.logs[lk] = c, l
	}
	return nil
}

func (f *fakeStore) Due(_ context.Context, user string, kind Kind, now time.Time, limit int) ([]string, error) {
	type kd struct {
		key string
		due time.Time
	}
	var out []kd
	for k, c := range f.cards {
		pre := user + "/" + string(kind) + "/"
		if len(k) > len(pre) && k[:len(pre)] == pre && !c.Due.After(now) {
			out = append(out, kd{k[len(pre):], c.Due})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].due.Before(out[j].due) })
	keys := []string{}
	for _, x := range out {
		if len(keys) < limit {
			keys = append(keys, x.key)
		}
	}
	return keys, nil
}

func TestRecordIsIdempotentPerAttempt(t *testing.T) {
	st := newFake()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Service{st: st, now: func() time.Time { return now }}
	ctx := context.Background()
	out := []Outcome{{KindFact, "f1", false}, {KindPattern, "p1", true}}
	if err := s.Record(ctx, "u", "a1", out); err != nil {
		t.Fatal(err)
	}
	before := st.cards["u/fact/f1"]
	if before.Reps != 1 || before.Lapses != 1 || !before.Due.After(now) {
		t.Fatalf("card = %+v", before)
	}
	now = now.Add(48 * time.Hour)
	if err := s.Record(ctx, "u", "a1", out); err != nil {
		t.Fatal(err)
	}
	if st.cards["u/fact/f1"] != before || len(st.logs) != 2 {
		t.Errorf("second Record changed state: %+v, %d logs", st.cards["u/fact/f1"], len(st.logs))
	}
	if err := s.Record(ctx, "u", "a2", out[:1]); err != nil {
		t.Fatal(err)
	}
	l := st.logs["a2/fact/f1"]
	if st.cards["u/fact/f1"].Reps != 2 || l.StabilityBefore != before.Stability || l.ElapsedDays != 2 || l.Grade != fsrs.Again {
		t.Errorf("log = %+v card = %+v", l, st.cards["u/fact/f1"])
	}
}

func TestDueAndCounts(t *testing.T) {
	st := newFake()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s := &Service{st: st, now: func() time.Time { return now }}
	ctx := context.Background()
	if err := s.Record(ctx, "u", "a1", []Outcome{{KindFact, "f1", false}, {KindFact, "f2", true}, {KindPattern, "p1", false}}); err != nil {
		t.Fatal(err)
	}
	if c, _ := s.Counts(ctx, "u"); c.DueCount != 0 {
		t.Errorf("counts before due = %+v", c)
	}
	now = now.Add(30 * 24 * time.Hour)
	c, err := s.Counts(ctx, "u")
	if err != nil {
		t.Fatal(err)
	}
	if c != (Counts{DueCount: 3, FactsDue: 2, PatternsDue: 1}) {
		t.Errorf("counts = %+v", c)
	}
	ids, _ := s.DueFacts(ctx, "u", 1)
	if len(ids) != 1 || ids[0] != "f1" { // Again has the shorter interval, so it is more overdue
		t.Errorf("due facts = %v", ids)
	}
	if other, _ := s.DueFacts(ctx, "other", 5); len(other) != 0 {
		t.Errorf("other user's due = %v", other)
	}
}
