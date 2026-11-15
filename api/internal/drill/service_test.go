package drill

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/question"
)

// These tests use the in-memory fakeStore, so they run without DATABASE_URL.

const fakeUser = "user-1"

func newFakeService() (*Service, *fakeStore) {
	st := newFakeStore()
	return New(st, catalog.All()), st
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
	svc := New(newFakeStore(), []gen.Family{brokenFamily{}})
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
