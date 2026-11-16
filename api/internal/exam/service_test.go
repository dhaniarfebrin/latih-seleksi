package exam

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"latihseleksi/api/internal/bank"
	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/spec"
)

// These tests use the in-memory fakeStore and an injected clock, so they run
// without DATABASE_URL.

type fakeEnv struct {
	svc   *Service
	st    *fakeStore
	clock time.Time
}

func newFakeEnv(t *testing.T) *fakeEnv {
	t.Helper()
	specs, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeEnv{st: newFakeStore(), clock: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)}
	f.svc = New(f.st, specs, NewGeneratedSource(catalog.All()))
	f.svc.now = func() time.Time { return f.clock }
	return f
}

func (f *fakeEnv) create(t *testing.T, user string) ExamView {
	t.Helper()
	v, err := f.svc.Create(context.Background(), user, "skd-2024", "umum", []string{"TIU"})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func (f *fakeEnv) snapshot(t *testing.T, id string, pos int) question.Snapshot {
	t.Helper()
	var sn question.Snapshot
	if err := json.Unmarshal(f.st.exams[id].items[pos-1].Snapshot, &sn); err != nil {
		t.Fatal(err)
	}
	return sn
}

func TestFakeCreate_startsTimerAndHidesKeys(t *testing.T) {
	f := newFakeEnv(t)
	v := f.create(t, "u1")
	if len(v.Items) != 35 || v.Status != StatusInProgress || !v.DeadlineAt.After(v.StartedAt) {
		t.Fatalf("view = %d items, status %s, %v..%v", len(v.Items), v.Status, v.StartedAt, v.DeadlineAt)
	}
	if !v.StartedAt.Equal(f.clock) {
		t.Errorf("started %v, want injected clock %v", v.StartedAt, f.clock)
	}
	raw, _ := json.Marshal(v)
	for _, banned := range []string{"answer_key", "steps", "misconception_id"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("open exam exposes %q", banned)
		}
	}
}

func TestFakeAnswer_highestSeqWinsThenSubmitScores(t *testing.T) {
	f := newFakeEnv(t)
	ctx := context.Background()
	v := f.create(t, "u1")
	key := f.snapshot(t, v.ID, 1).AnswerKey
	wrong := "A"
	if key == "A" {
		wrong = "B"
	}
	r, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: wrong, TimeMS: 100, Seq: 2})
	if err != nil || r.Answer != wrong {
		t.Fatalf("first = %+v, %v", r, err)
	}
	r, _ = f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: key, TimeMS: 50, Seq: 1})
	if r.Answer != wrong || r.Seq != 2 {
		t.Errorf("stale seq replaced the answer: %+v", r)
	}
	r, _ = f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: key, TimeMS: 300, Flagged: true, Seq: 3})
	if r.Answer != key || !r.Flagged {
		t.Errorf("newer seq not stored: %+v", r)
	}
	res, err := f.svc.Submit(ctx, "u1", v.ID)
	if err != nil || res.Status != StatusSubmitted || res.Sections[0].Score != 5 || res.Sections[0].Answered != 1 {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 2, Answer: "A", Seq: 1}); !errors.Is(err, ErrClosed) {
		t.Errorf("answer after submit err = %v, want ErrClosed", err)
	}
}

func TestFakeLoad_expiresAfterDeadlinePlusGrace(t *testing.T) {
	f := newFakeEnv(t)
	ctx := context.Background()
	v := f.create(t, "u1")
	f.clock = v.DeadlineAt.Add(grace - time.Second)
	if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: "A", Seq: 1}); err != nil {
		t.Fatalf("answer inside grace: %v", err)
	}
	f.clock = v.DeadlineAt.Add(grace + time.Second)
	g, err := f.svc.Get(ctx, "u1", v.ID)
	if err != nil || g.Status != StatusExpired || g.SubmittedAt == nil || !g.SubmittedAt.Equal(v.DeadlineAt) {
		t.Fatalf("get = %+v, %v", g, err)
	}
	if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: "B", Seq: 2}); !errors.Is(err, ErrClosed) {
		t.Errorf("answer after expiry err = %v, want ErrClosed", err)
	}
}

func TestFakeAccess_ownershipAndErrors(t *testing.T) {
	f := newFakeEnv(t)
	ctx := context.Background()
	v := f.create(t, "u1")
	if _, err := f.svc.Get(ctx, "u2", v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("other user err = %v, want ErrNotFound", err)
	}
	if _, err := f.svc.Get(ctx, "", v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("anonymous err = %v, want ErrNotFound", err)
	}
	if _, err := f.svc.Result(ctx, "u1", v.ID); !errors.Is(err, ErrInProgress) {
		t.Errorf("result while open err = %v, want ErrInProgress", err)
	}
	cases := []struct {
		name string
		in   AnswerInput
		want error
	}{
		{"bad letter", AnswerInput{Pos: 1, Answer: "Z"}, ErrInvalid},
		{"negative time", AnswerInput{Pos: 1, Answer: "A", TimeMS: -1}, ErrInvalid},
		{"negative seq", AnswerInput{Pos: 1, Answer: "A", Seq: -1}, ErrInvalid},
		{"pos out of range", AnswerInput{Pos: 99, Answer: "A"}, ErrNotFound},
	}
	for _, c := range cases {
		if _, err := f.svc.Answer(ctx, "u1", v.ID, c.in); !errors.Is(err, c.want) {
			t.Errorf("%s: err %v, want %v", c.name, err, c.want)
		}
	}
}

func TestFakeList_newestFirstPerUser(t *testing.T) {
	f := newFakeEnv(t)
	ctx := context.Background()
	a := f.create(t, "u1")
	f.create(t, "u2")
	b := f.create(t, "u1")
	if _, err := f.svc.Submit(ctx, "u1", a.ID); err != nil {
		t.Fatal(err)
	}
	list, err := f.svc.List(ctx, "u1", 10)
	if err != nil || len(list) != 2 || list[0].ID != b.ID || list[1].ID != a.ID {
		t.Fatalf("list = %+v, %v", list, err)
	}
	if list[0].Result != nil || list[1].Result == nil {
		t.Errorf("results: open %v, submitted %v", list[0].Result, list[1].Result)
	}
	if _, err := f.svc.List(ctx, "u1", 0); !errors.Is(err, ErrInvalid) {
		t.Errorf("limit 0 err = %v, want ErrInvalid", err)
	}
}

// A source failure must not be reported as bad input, and must not reach the store.
func TestFakeCreate_sourceFailureLeavesNoExam(t *testing.T) {
	f := newFakeEnv(t)
	specs, _ := spec.Load()
	svc := New(f.st, specs, failingSource{})
	_, err := svc.Create(context.Background(), "u1", "skd-2024", "umum", []string{"TIU"})
	if err == nil || errors.Is(err, ErrInvalid) {
		t.Fatalf("err = %v, want a non-ErrInvalid error", err)
	}
	if len(f.st.exams) != 0 {
		t.Errorf("%d exams stored after a failed create", len(f.st.exams))
	}
}

// weightedSource serves 45 TKP items whose option A..E weigh 1..5; the key is
// the weight-5 option E. badWeight, when set, gives option A that weight.
type weightedSource struct{ badWeight int }

func (weightedSource) Sections() []string { return []string{"TKP"} }
func (w weightedSource) Pick(_ string, n int, _ int64) ([]gen.Draft, error) {
	out := make([]gen.Draft, 0, n)
	for i := 0; i < n; i++ {
		d := gen.Draft{PatternID: "tkp-test", Stem: "s", AnswerKey: "E"}
		for k, key := range []string{"A", "B", "C", "D", "E"} {
			wt := k + 1
			if key == "A" && w.badWeight != 0 {
				wt = w.badWeight
			}
			d.Options = append(d.Options, gen.Option{Key: key, Text: key, Weight: wt})
		}
		out = append(out, d)
	}
	return out, nil
}

func newWeightedEnv(t *testing.T, src ItemSource) *fakeEnv {
	t.Helper()
	f := newFakeEnv(t)
	specs, _ := spec.Load()
	f.svc = New(f.st, specs, src)
	f.svc.now = func() time.Time { return f.clock }
	return f
}

func TestFakeWeighted_availableAndScoresSumOfWeights(t *testing.T) {
	f := newWeightedEnv(t, weightedSource{})
	ctx := context.Background()
	for _, sp := range f.svc.Specs() {
		for _, sec := range sp.Sections {
			if sec.ID == "TKP" && !sec.Available {
				t.Errorf("%s TKP not available with a TKP source", sp.ID)
			}
		}
	}
	v, err := f.svc.Create(ctx, "u1", "skd-2024", "umum", []string{"TKP"})
	if err != nil || len(v.Items) != 45 {
		t.Fatalf("create = %d items, %v", len(v.Items), err)
	}
	// 3 x E (5), 2 x B (2), 1 x A (1), rest blank (blank = 0 in the spec).
	answers := map[int]string{1: "E", 2: "E", 3: "E", 4: "B", 5: "B", 6: "A"}
	for pos, a := range answers {
		if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: pos, Answer: a, Seq: 1}); err != nil {
			t.Fatal(err)
		}
	}
	res, err := f.svc.Submit(ctx, "u1", v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got := res.Sections[0]; got.Score != 15+4+1 || got.Answered != 6 || got.Max != 225 {
		t.Errorf("section = %+v, want score 20 answered 6 max 225", got)
	}
	g, err := f.svc.Get(ctx, "u1", v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if r := g.Items[3].Response; r == nil || r.Weight == nil || *r.Weight != 2 {
		t.Errorf("closed TKP response weight = %+v, want 2", r)
	}
	if r := g.Items[6].Response; r != nil {
		t.Errorf("unanswered item has response %+v", r)
	}
}

func TestFakeWeighted_openExamHidesWeightAndRightWrongHasNone(t *testing.T) {
	f := newWeightedEnv(t, weightedSource{})
	ctx := context.Background()
	v, _ := f.svc.Create(ctx, "u1", "skd-2024", "umum", []string{"TKP"})
	if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: "E", Seq: 1}); err != nil {
		t.Fatal(err)
	}
	g, _ := f.svc.Get(ctx, "u1", v.ID)
	raw, _ := json.Marshal(g)
	if strings.Contains(string(raw), "weight") {
		t.Errorf("open exam exposes weight: %s", raw)
	}
	tiu := newFakeEnv(t)
	tv := tiu.create(t, "u1")
	_, _ = tiu.svc.Answer(ctx, "u1", tv.ID, AnswerInput{Pos: 1, Answer: "A", Seq: 1})
	_, _ = tiu.svc.Submit(ctx, "u1", tv.ID)
	tg, _ := tiu.svc.Get(ctx, "u1", tv.ID)
	if tg.Items[0].Response.Weight != nil {
		t.Errorf("TIU response has weight %v", *tg.Items[0].Response.Weight)
	}
}

func TestFakeWeighted_weightOutsideRangeIsError(t *testing.T) {
	f := newWeightedEnv(t, weightedSource{badWeight: 9})
	ctx := context.Background()
	v, err := f.svc.Create(ctx, "u1", "skd-2024", "umum", []string{"TKP"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: 1, Answer: "A", Seq: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.Submit(ctx, "u1", v.ID); err == nil || !strings.Contains(err.Error(), "weight 9") {
		t.Errorf("Submit err = %v, want a weight 9 error", err)
	}
}

func TestFakeFullSpec_realBankAndGeneratedMake110(t *testing.T) {
	b, err := bank.Load()
	if err != nil {
		t.Fatal(err)
	}
	f := newFakeEnv(t)
	specs, _ := spec.Load()
	f.svc = New(f.st, specs, NewGeneratedSource(catalog.All()), b)
	f.svc.now = func() time.Time { return f.clock }
	ctx := context.Background()
	v, err := f.svc.Create(ctx, "u1", "skd-2024", "umum", []string{"TWK", "TIU", "TKP"})
	if err != nil || len(v.Items) != 110 {
		t.Fatalf("create = %d items, %v", len(v.Items), err)
	}
	for pos := 1; pos <= 110; pos++ {
		if _, err := f.svc.Answer(ctx, "u1", v.ID, AnswerInput{Pos: pos, Answer: "A", Seq: 1}); err != nil {
			t.Fatalf("answer %d: %v", pos, err)
		}
	}
	res, err := f.svc.Submit(ctx, "u1", v.ID)
	if err != nil || res.Partial || len(res.Sections) != 3 {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if tkp := res.Sections[2]; tkp.ID != "TKP" || tkp.Score < 45 || tkp.Score > 225 {
		t.Errorf("TKP section = %+v, want 45..225 for 45 answered items", tkp)
	}
}
