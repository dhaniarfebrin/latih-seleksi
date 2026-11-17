package drill

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/review"
)

// These tests use the in-memory fakeStore, so they run without DATABASE_URL.

const fakeUser = "user-1"

func newFakeService() (*Service, *fakeStore) {
	st := newFakeStore()
	return New(st, catalog.All(), nil, nil), st
}

func fakeSnapshot(t *testing.T, st *fakeStore, id string, pos int) question.Snapshot {
	t.Helper()
	var sn question.Snapshot
	if err := json.Unmarshal(st.attempts[id].items[pos-1].Snapshot, &sn); err != nil {
		t.Fatal(err)
	}
	return sn
}

func TestFakeCreate_freezesItemsAndHidesKey(t *testing.T) {
	svc, st := newFakeService()
	v, err := svc.Create(context.Background(), fakeUser, persentase.PatternID, 5, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(v.Items) != 5 || v.N != 5 || v.Mode != ModeDrill {
		t.Fatalf("view = %d items, n %d, mode %s", len(v.Items), v.N, v.Mode)
	}
	raw, _ := json.Marshal(v)
	for _, banned := range []string{"answer_key", "steps", "misconception_id"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("create response contains %q", banned)
		}
	}
	if fakeSnapshot(t, st, v.ID, 1).AnswerKey == "" {
		t.Error("stored snapshot has no answer key")
	}
}

func TestFakeCreate_invalidInput(t *testing.T) {
	svc, _ := newFakeService()
	cases := []struct {
		pattern string
		n, band int
	}{
		{"nope", 3, 2}, {persentase.PatternID, 0, 2}, {persentase.PatternID, 21, 2},
		{persentase.PatternID, 3, 0}, {persentase.PatternID, 3, 4},
	}
	for _, c := range cases {
		if _, err := svc.Create(context.Background(), fakeUser, c.pattern, c.n, c.band); !errors.Is(err, ErrInvalid) {
			t.Errorf("Create(%q,%d,%d) err = %v, want ErrInvalid", c.pattern, c.n, c.band, err)
		}
	}
}

type brokenFamily struct{}

func (brokenFamily) ID() string                        { return "broken" }
func (brokenFamily) Label() string                     { return "Broken" }
func (brokenFamily) Misconceptions() map[string]string { return nil }
func (brokenFamily) Generate(int64, gen.Band) (gen.Draft, error) {
	return gen.Draft{}, errFamilyBroken
}

// A generator failure is a server fault, so it must not look like bad input.
func TestFakeCreate_generatorFailureIsNotInvalid(t *testing.T) {
	svc := New(newFakeStore(), []gen.Family{brokenFamily{}}, nil, nil)
	_, err := svc.Create(context.Background(), fakeUser, "broken", 3, 1)
	if !errors.Is(err, errFamilyBroken) || errors.Is(err, ErrInvalid) {
		t.Fatalf("err = %v, want errFamilyBroken and not ErrInvalid", err)
	}
}

func TestFakeAnswer_revealsIsIdempotentAndScores(t *testing.T) {
	svc, st := newFakeService()
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, persentase.PatternID, 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	sn := fakeSnapshot(t, st, v.ID, 1)
	first, err := svc.Answer(ctx, fakeUser, v.ID, 1, sn.AnswerKey, 1200)
	if err != nil || !first.Correct || first.AnswerKey != sn.AnswerKey || len(first.Steps) == 0 {
		t.Fatalf("first = %+v, %v", first, err)
	}
	wrong := "A"
	if sn.AnswerKey == "A" {
		wrong = "B"
	}
	again, err := svc.Answer(ctx, fakeUser, v.ID, 1, wrong, 9999)
	if err != nil || !again.Correct || again.TimeMS != 1200 {
		t.Fatalf("repeat changed the stored answer: %+v, %v", again, err)
	}
	res, err := svc.Submit(ctx, fakeUser, v.ID)
	if err != nil || res.Score != 1 || res.Answered != 1 || res.N != 3 || res.AvgTimeMS != 1200 {
		t.Fatalf("result = %+v, %v", res, err)
	}
	if _, err := svc.Answer(ctx, fakeUser, v.ID, 2, "A", 10); !errors.Is(err, ErrSubmitted) {
		t.Errorf("answer after submit err = %v, want ErrSubmitted", err)
	}
}

func TestFakeAnswer_invalidAndMissing(t *testing.T) {
	svc, _ := newFakeService()
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, persentase.PatternID, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, ans := range []string{"", "Z", "a", "AB"} {
		if _, err := svc.Answer(ctx, fakeUser, v.ID, 1, ans, 100); !errors.Is(err, ErrInvalid) {
			t.Errorf("answer %q err = %v, want ErrInvalid", ans, err)
		}
	}
	if _, err := svc.Answer(ctx, fakeUser, v.ID, 1, "A", -1); !errors.Is(err, ErrInvalid) {
		t.Errorf("negative time err = %v, want ErrInvalid", err)
	}
	if _, err := svc.Answer(ctx, fakeUser, v.ID, 9, "A", 100); !errors.Is(err, ErrNotFound) {
		t.Errorf("pos 9 err = %v, want ErrNotFound", err)
	}
	if _, err := svc.Get(ctx, fakeUser, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get err = %v, want ErrNotFound", err)
	}
	if _, err := svc.Submit(ctx, fakeUser, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Submit err = %v, want ErrNotFound", err)
	}
}

func TestFakeOwnership_otherUserSeesNotFound(t *testing.T) {
	svc, _ := newFakeService()
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, persentase.PatternID, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	const other = "user-2"
	if _, err := svc.Get(ctx, other, v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get err = %v, want ErrNotFound", err)
	}
	if _, err := svc.Answer(ctx, other, v.ID, 1, "A", 10); !errors.Is(err, ErrNotFound) {
		t.Errorf("Answer err = %v, want ErrNotFound", err)
	}
	if _, err := svc.Submit(ctx, other, v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Submit err = %v, want ErrNotFound", err)
	}
	if got, err := svc.Get(ctx, fakeUser, v.ID); err != nil || got.SubmittedAt != nil {
		t.Errorf("owner view changed: %+v, %v", got, err)
	}
}

func TestFakeReport(t *testing.T) {
	svc, st := newFakeService()
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, persentase.PatternID, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Report(ctx, fakeUser, v.ID, 1, "typo", ""); !errors.Is(err, ErrInvalid) {
		t.Errorf("unanswered err = %v, want ErrInvalid", err)
	}
	if _, err := svc.Answer(ctx, fakeUser, v.ID, 1, "A", 100); err != nil {
		t.Fatal(err)
	}
	if err := svc.Report(ctx, fakeUser, v.ID, 1, "nope", ""); !errors.Is(err, ErrInvalid) {
		t.Errorf("bad reason err = %v, want ErrInvalid", err)
	}
	if err := svc.Report(ctx, fakeUser, v.ID, 1, "typo", strings.Repeat("x", question.MaxReportNote+1)); !errors.Is(err, ErrInvalid) {
		t.Errorf("long note err = %v, want ErrInvalid", err)
	}
	if err := svc.Report(ctx, "user-2", v.ID, 1, "typo", ""); !errors.Is(err, ErrNotFound) {
		t.Errorf("other user err = %v, want ErrNotFound", err)
	}
	if err := svc.Report(ctx, fakeUser, v.ID, 9, "typo", ""); !errors.Is(err, ErrNotFound) {
		t.Errorf("pos 9 err = %v, want ErrNotFound", err)
	}
	if len(st.reports) != 0 {
		t.Fatalf("rejected reports were stored: %v", st.reports)
	}
	if err := svc.Report(ctx, fakeUser, v.ID, 1, "typo", "salah ketik"); err != nil {
		t.Fatal(err)
	}
	if got := st.reports[v.ID+"/1"]; got != (fakeReport{"typo", "salah ketik"}) {
		t.Errorf("stored = %+v", got)
	}
}

// fakeFacts serves two TWK drafts for one pattern.
type fakeFacts struct{}

func (fakeFacts) PatternIDs() []string { return []string{"TWK.x"} }
func (fakeFacts) PatternLabel(id string) (string, bool) {
	return "TWK: x", id == "TWK.x"
}
func (fakeFacts) Misconceptions() map[string]string { return map[string]string{"distraktor": "d"} }
func (fakeFacts) Draw(_ string, n int, seed int64) ([]gen.Draft, error) {
	var out []gen.Draft
	for i := 0; i < min(n, 2); i++ {
		out = append(out, gen.Draft{
			Stem: "q" + string(rune('0'+i)), Seed: seed + int64(i), FactID: "fact-" + string(rune('0'+i)), AnswerKey: "A",
			Options: []gen.Option{{Key: "A", Text: "ya"}, {Key: "B", Text: "tidak", MisconceptionID: "distraktor"}},
			Steps:   []string{"Pasal 1: \"ya\""},
		})
	}
	return out, nil
}

func (fakeFacts) ByID(ids []string, seed int64) []gen.Draft {
	var out []gen.Draft
	for i, id := range ids {
		if id == "gone" {
			continue
		}
		out = append(out, gen.Draft{
			PatternID: "TWK.x", Stem: "r-" + id, Seed: seed + int64(i), FactID: id, AnswerKey: "A",
			Options: []gen.Option{{Key: "A", Text: "ya"}, {Key: "B", Text: "tidak", MisconceptionID: "distraktor"}},
			Steps:   []string{"Pasal 1: \"ya\""},
		})
	}
	return out
}

// fakeReviewer keeps what Record received and serves a fixed due list.
type fakeReviewer struct {
	due      []string
	recorded map[string][]review.Outcome
	calls    int
}

func (f *fakeReviewer) Record(_ context.Context, _, attemptID string, o []review.Outcome) error {
	if f.recorded == nil {
		f.recorded = map[string][]review.Outcome{}
	}
	f.recorded[attemptID] = o
	f.calls++
	return nil
}

func (f *fakeReviewer) DueFacts(_ context.Context, _ string, limit int) ([]string, error) {
	return f.due[:min(limit, len(f.due))], nil
}

func TestSubmitRecordsFactOutcomes(t *testing.T) {
	st, rev := newFakeStore(), &fakeReviewer{}
	svc := New(st, catalog.All(), fakeFacts{}, rev)
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, "TWK.x", 8, 2)
	if err != nil {
		t.Fatal(err)
	}
	svc.Answer(ctx, fakeUser, v.ID, 1, "A", 10) //nolint:errcheck
	svc.Answer(ctx, fakeUser, v.ID, 2, "B", 10) //nolint:errcheck
	for i := 0; i < 2; i++ {                    // repeating Submit sends the same outcomes again
		if _, err := svc.Submit(ctx, fakeUser, v.ID); err != nil {
			t.Fatal(err)
		}
	}
	want := []review.Outcome{{Kind: review.KindFact, Key: "fact-0", Correct: true}, {Kind: review.KindFact, Key: "fact-1", Correct: false}}
	if !reflect.DeepEqual(rev.recorded[v.ID], want) || rev.calls != 2 {
		t.Errorf("recorded = %+v after %d calls", rev.recorded[v.ID], rev.calls)
	}
}

func TestSubmitRecordsPatternOutcome(t *testing.T) {
	cases := []struct {
		name    string
		answers []string // "ok" or "bad" per position
		want    bool
		none    bool
	}{
		{"4 of 5 is good", []string{"ok", "ok", "ok", "ok", "bad"}, true, false},
		{"3 of 5 is again", []string{"ok", "ok", "ok", "bad", "bad"}, false, false},
		{"no answers schedules nothing", nil, false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			st, rev := newFakeStore(), &fakeReviewer{}
			svc := New(st, catalog.All(), nil, rev)
			ctx := context.Background()
			v, err := svc.Create(ctx, fakeUser, persentase.PatternID, 5, 2)
			if err != nil {
				t.Fatal(err)
			}
			for i, kind := range c.answers {
				key := fakeSnapshot(t, st, v.ID, i+1).AnswerKey
				if kind == "bad" {
					key = map[string]string{"A": "B"}[key]
					if key == "" {
						key = "A"
					}
				}
				if _, err := svc.Answer(ctx, fakeUser, v.ID, i+1, key, 10); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := svc.Submit(ctx, fakeUser, v.ID); err != nil {
				t.Fatal(err)
			}
			got := rev.recorded[v.ID]
			if c.none {
				if len(got) != 0 {
					t.Errorf("recorded = %+v", got)
				}
				return
			}
			if len(got) != 1 || got[0] != (review.Outcome{Kind: review.KindPattern, Key: persentase.PatternID, Correct: c.want}) {
				t.Errorf("recorded = %+v", got)
			}
		})
	}
}

func TestCreateReviewDrill(t *testing.T) {
	st, rev := newFakeStore(), &fakeReviewer{due: []string{"fact-9", "gone", "fact-3"}}
	svc := New(st, catalog.All(), fakeFacts{}, rev)
	ctx := context.Background()
	v, err := svc.Create(ctx, fakeUser, ReviewPatternID, 8, 2)
	if err != nil {
		t.Fatal(err)
	}
	if v.PatternID != ReviewPatternID || v.N != 2 || len(v.Items) != 2 {
		t.Fatalf("view = pattern %s n %d items %d", v.PatternID, v.N, len(v.Items))
	}
	items := st.attempts[v.ID].items
	if items[0].PatternID != "TWK.x" || items[0].Section != "TWK" {
		t.Errorf("item = %+v, want its own pattern and section", items[0])
	}
	if sn := fakeSnapshot(t, st, v.ID, 1); sn.FactID != "fact-9" {
		t.Errorf("first fact = %q, want the most overdue", sn.FactID)
	}
	// Answers in a review drill review the fact only, not the ULANG pattern.
	svc.Answer(ctx, fakeUser, v.ID, 1, "A", 10) //nolint:errcheck
	if _, err := svc.Submit(ctx, fakeUser, v.ID); err != nil {
		t.Fatal(err)
	}
	if got := rev.recorded[v.ID]; len(got) != 1 || got[0].Kind != review.KindFact {
		t.Errorf("recorded = %+v", got)
	}

	rev.due = nil
	if _, err := svc.Create(ctx, fakeUser, ReviewPatternID, 8, 2); !errors.Is(err, ErrInvalid) {
		t.Errorf("nothing due err = %v", err)
	}
	if _, err := New(st, catalog.All(), fakeFacts{}, nil).Create(ctx, fakeUser, ReviewPatternID, 8, 2); !errors.Is(err, ErrInvalid) {
		t.Errorf("without reviewer err = %v", err)
	}
}

// A stale due id (no longer reviewed in the bank) among the most overdue must
// not hide valid ones: with n = 1 the first id is dropped and fact-3 is served.
func TestCreateReviewSkipsStaleIDs(t *testing.T) {
	st, rev := newFakeStore(), &fakeReviewer{due: []string{"gone", "fact-3", "fact-4"}}
	svc := New(st, catalog.All(), fakeFacts{}, rev)
	v, err := svc.Create(context.Background(), fakeUser, ReviewPatternID, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if v.N != 1 || len(v.Items) != 1 {
		t.Fatalf("n %d items %d, want 1", v.N, len(v.Items))
	}
	if sn := fakeSnapshot(t, st, v.ID, 1); sn.FactID != "fact-3" {
		t.Errorf("fact = %q, want fact-3", sn.FactID)
	}
}

func TestFakeFactsPatternFlow(t *testing.T) {
	st := newFakeStore()
	svc := New(st, catalog.All(), fakeFacts{}, nil)
	var last Pattern
	for _, p := range svc.Patterns() {
		last = p
	}
	if last.ID != "TWK.x" || last.Misconceptions["distraktor"] == "" {
		t.Fatalf("last pattern = %+v", last)
	}
	v, err := svc.Create(context.Background(), fakeUser, "TWK.x", 8, 2)
	if err != nil {
		t.Fatal(err)
	}
	if v.N != 2 || len(v.Items) != 2 {
		t.Fatalf("n = %d, items = %d, want the pool size 2", v.N, len(v.Items))
	}
	if sn := fakeSnapshot(t, st, v.ID, 1); sn.FactID != "fact-0" {
		t.Errorf("fact id = %q", sn.FactID)
	}
	r, err := svc.Answer(context.Background(), fakeUser, v.ID, 1, "B", 1000)
	if err != nil || r.Correct || r.AnswerKey != "A" || r.MisconceptionID != "distraktor" || len(r.Steps) != 1 {
		t.Fatalf("answer = %+v, %v", r, err)
	}
	if _, err := svc.Create(context.Background(), fakeUser, "TWK.nope", 3, 2); !errors.Is(err, ErrInvalid) {
		t.Errorf("unknown pattern err = %v", err)
	}
}
