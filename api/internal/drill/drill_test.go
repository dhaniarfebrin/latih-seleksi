package drill

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/kerjasama"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/store"
)

func newService(t *testing.T) (*Service, *store.Store, string) {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	st, err := store.Open(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	uid, err := st.UpsertUser(context.Background(), "drill-test-"+t.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.DeleteUser(context.Background(), uid) })
	return New(st, catalog.All()), st, uid
}

// frozenSnapshot reads the stored question.Snapshot of pos through the store.
func frozenSnapshot(t *testing.T, st *store.Store, uid, id string, pos int) question.Snapshot {
	t.Helper()
	_, items, _, err := st.GetAttempt(context.Background(), uid, id)
	if err != nil {
		t.Fatal(err)
	}
	var sn question.Snapshot
	if err := json.Unmarshal(items[pos-1].Snapshot, &sn); err != nil {
		t.Fatal(err)
	}
	return sn
}

func TestPatternsListAllFamiliesWithLabels(t *testing.T) {
	got := New(nil, catalog.All()).Patterns() // Patterns does not touch the store
	var ids []string
	for _, p := range got {
		ids = append(ids, p.ID)
		if p.Label == "" || len(p.Misconceptions) == 0 {
			t.Errorf("pattern %s: label %q, %d misconceptions", p.ID, p.Label, len(p.Misconceptions))
		}
		for id, text := range p.Misconceptions {
			if text == "" {
				t.Errorf("pattern %s: empty label for %q", p.ID, id)
			}
		}
	}
	var want []string
	for _, f := range catalog.All() {
		want = append(want, f.ID())
	}
	if strings.Join(ids, ",") != strings.Join(want, ",") {
		t.Fatalf("pattern ids = %v, want %v", ids, want)
	}
}

// kerja-sama band 1 has only a few hundred distinct stems, so 20 random items
// would repeat one about half the time without the de-duplication in Create.
func TestCreateAvoidsRepeatedStems(t *testing.T) {
	svc, _, uid := newService(t)
	v, err := svc.Create(context.Background(), uid, kerjasama.PatternID, 20, 1)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, it := range v.Items {
		if seen[it.Stem] {
			t.Errorf("repeated stem at pos %d: %s", it.Pos, it.Stem)
		}
		seen[it.Stem] = true
	}
}

func TestCreateWorksForEveryPattern(t *testing.T) {
	svc, _, uid := newService(t)
	for _, p := range svc.Patterns() {
		v, err := svc.Create(context.Background(), uid, p.ID, 3, 3)
		if err != nil {
			t.Fatalf("Create(%s): %v", p.ID, err)
		}
		if v.PatternID != p.ID || len(v.Items) != 3 {
			t.Errorf("Create(%s) = pattern %s, %d items", p.ID, v.PatternID, len(v.Items))
		}
	}
}

func TestCreateHidesKey(t *testing.T) {
	svc, _, uid := newService(t)
	v, err := svc.Create(context.Background(), uid, persentase.PatternID, 8, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(v.Items) != 8 {
		t.Fatalf("items = %d, want 8", len(v.Items))
	}
	raw, _ := json.Marshal(v)
	for _, banned := range []string{"answer_key", "steps", "misconception_id"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("create response contains %q", banned)
		}
	}
	got, err := svc.Get(context.Background(), uid, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ = json.Marshal(got)
	for _, banned := range []string{"answer_key", "steps", "misconception_id"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("get response contains %q", banned)
		}
	}
}

func TestAnswerRevealsAndIsIdempotent(t *testing.T) {
	svc, st, uid := newService(t)
	ctx := context.Background()
	v, err := svc.Create(ctx, uid, persentase.PatternID, 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	sn := frozenSnapshot(t, st, uid, v.ID, 1)

	first, err := svc.Answer(ctx, uid, v.ID, 1, sn.AnswerKey, 1200)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Correct || first.AnswerKey != sn.AnswerKey || len(first.Steps) == 0 || first.MisconceptionID != "" {
		t.Fatalf("unexpected correct reveal: %+v", first)
	}

	// A different answer and time for the same position must not change the result.
	wrong := "A"
	if sn.AnswerKey == "A" {
		wrong = "B"
	}
	again, err := svc.Answer(ctx, uid, v.ID, 1, wrong, 9999)
	if err != nil {
		t.Fatal(err)
	}
	if again.Answer != first.Answer || again.TimeMS != first.TimeMS || !again.Correct {
		t.Fatalf("second PUT changed result: first=%+v again=%+v", first, again)
	}

	// A wrong answer at pos 2 carries its misconception id.
	sn2 := frozenSnapshot(t, st, uid, v.ID, 2)
	var wrongOpt question.SnapshotOption
	for _, o := range sn2.Options {
		if o.Key != sn2.AnswerKey {
			wrongOpt = o
			break
		}
	}
	r2, err := svc.Answer(ctx, uid, v.ID, 2, wrongOpt.Key, 3000)
	if err != nil {
		t.Fatal(err)
	}
	if r2.Correct || r2.MisconceptionID == "" || r2.MisconceptionID != wrongOpt.MisconceptionID {
		t.Fatalf("unexpected wrong reveal: %+v", r2)
	}

	got, err := svc.Get(ctx, uid, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Items[0].Response == nil || got.Items[1].Response == nil || got.Items[2].Response != nil {
		t.Fatalf("response visibility wrong after answering pos 1 and 2")
	}
}

func TestSubmitScores(t *testing.T) {
	svc, st, uid := newService(t)
	ctx := context.Background()
	v, err := svc.Create(ctx, uid, persentase.PatternID, 4, 2)
	if err != nil {
		t.Fatal(err)
	}
	wantMis := map[string]int{}
	wantScore := 0
	// pos 1,2 correct (1000, 2000 ms); pos 3 wrong (3000 ms); pos 4 unanswered.
	for pos, ms := range map[int]int{1: 1000, 2: 2000, 3: 3000} {
		sn := frozenSnapshot(t, st, uid, v.ID, pos)
		ans := sn.AnswerKey
		if pos == 3 {
			for _, o := range sn.Options {
				if o.Key != sn.AnswerKey {
					ans = o.Key
					wantMis[o.MisconceptionID]++
					break
				}
			}
		} else {
			wantScore++
		}
		if _, err := svc.Answer(ctx, uid, v.ID, pos, ans, ms); err != nil {
			t.Fatal(err)
		}
	}
	res, err := svc.Submit(ctx, uid, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if res.Score != wantScore || res.N != 4 || res.Answered != 3 || res.AvgTimeMS != 2000 {
		t.Fatalf("result = %+v", res)
	}
	if len(res.Misconception) != len(wantMis) {
		t.Fatalf("misconceptions = %v, want %v", res.Misconception, wantMis)
	}
	for k, n := range wantMis {
		if res.Misconception[k] != n {
			t.Fatalf("misconceptions = %v, want %v", res.Misconception, wantMis)
		}
	}
	again, err := svc.Submit(ctx, uid, v.ID)
	if err != nil || !again.SubmittedAt.Equal(res.SubmittedAt) || again.Score != res.Score {
		t.Fatalf("second submit differs: %+v err=%v", again, err)
	}
	// Unanswered position is closed after submit.
	if _, err := svc.Answer(ctx, uid, v.ID, 4, "A", 10); !errors.Is(err, ErrSubmitted) {
		t.Fatalf("answer after submit err = %v, want ErrSubmitted", err)
	}
}

func TestValidation(t *testing.T) {
	svc, _, uid := newService(t)
	ctx := context.Background()
	for _, c := range []struct {
		pattern string
		n, band int
	}{
		{persentase.PatternID, 0, 2}, {persentase.PatternID, 21, 2},
		{persentase.PatternID, 8, 0}, {persentase.PatternID, 8, 4},
		{"nope", 8, 2},
	} {
		if _, err := svc.Create(ctx, uid, c.pattern, c.n, c.band); !errors.Is(err, ErrInvalid) {
			t.Errorf("Create(%q,%d,%d) err = %v, want ErrInvalid", c.pattern, c.n, c.band, err)
		}
	}
	v, err := svc.Create(ctx, uid, persentase.PatternID, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, ans := range []string{"", "F", "a", "AB"} {
		if _, err := svc.Answer(ctx, uid, v.ID, 1, ans, 100); !errors.Is(err, ErrInvalid) {
			t.Errorf("answer %q err = %v, want ErrInvalid", ans, err)
		}
	}
	if _, err := svc.Answer(ctx, uid, v.ID, 1, "A", -1); !errors.Is(err, ErrInvalid) {
		t.Errorf("negative time err = %v, want ErrInvalid", err)
	}
	for _, pos := range []int{0, 3, -1} {
		if _, err := svc.Answer(ctx, uid, v.ID, pos, "A", 100); !errors.Is(err, ErrNotFound) {
			t.Errorf("pos %d err = %v, want ErrNotFound", pos, err)
		}
	}
	for _, id := range []string{"not-a-uuid", "00000000-0000-0000-0000-000000000000"} {
		if _, err := svc.Get(ctx, uid, id); !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(%q) err = %v, want ErrNotFound", id, err)
		}
		if _, err := svc.Submit(ctx, uid, id); !errors.Is(err, ErrNotFound) {
			t.Errorf("Submit(%q) err = %v, want ErrNotFound", id, err)
		}
	}
}
