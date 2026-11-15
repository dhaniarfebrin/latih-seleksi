package exam

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/scoring"
	"latihseleksi/api/internal/spec"
	"latihseleksi/api/internal/store"
)

func TestGeneratedSourceDistinctStemsAndDeterministic(t *testing.T) {
	src := NewGeneratedSource(catalog.All())
	for seed := int64(1); seed <= 50; seed++ {
		a, err := src.Pick("TIU", 35, seed)
		if err != nil {
			t.Fatal(err)
		}
		if len(a) != 35 {
			t.Fatalf("seed %d: %d items", seed, len(a))
		}
		seen := map[string]bool{}
		for _, d := range a {
			if seen[d.Stem] {
				t.Fatalf("seed %d: repeated stem %q", seed, d.Stem)
			}
			seen[d.Stem] = true
		}
		b, _ := src.Pick("TIU", 35, seed)
		for i := range a {
			if a[i].Stem != b[i].Stem || a[i].AnswerKey != b[i].AnswerKey {
				t.Fatalf("seed %d: item %d differs between runs", seed, i)
			}
		}
	}
	if _, err := src.Pick("TWK", 1, 1); err == nil {
		t.Error("TWK accepted by the generated source")
	}
}

type env struct {
	user  string
	svc   *Service
	st    *store.Store
	clock time.Time
}

func newEnv(t *testing.T) *env {
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
	specs, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	e := &env{st: st, clock: time.Now().UTC().Truncate(time.Millisecond)}
	e.user = e.newUser(t)
	e.svc = New(st, specs, NewGeneratedSource(catalog.All()))
	e.svc.now = func() time.Time { return e.clock }
	return e
}

// newUser makes a fresh user so tests do not see each other's exams.
func (e *env) newUser(t *testing.T) string {
	t.Helper()
	var b [8]byte
	_, _ = rand.Read(b[:])
	id, err := e.st.UpsertUser(context.Background(), "test-"+hex.EncodeToString(b[:]))
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func (e *env) create(t *testing.T) ExamView {
	t.Helper()
	v, err := e.svc.Create(context.Background(), e.user, "skd-2024", "umum", []string{"TIU"})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// keys reads the frozen answer keys straight from the store.
func (e *env) keys(t *testing.T, id string) map[int]snapshot {
	t.Helper()
	_, items, _, err := e.st.GetExam(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	out := map[int]snapshot{}
	for _, it := range items {
		var sn snapshot
		if err := json.Unmarshal(it.Snapshot, &sn); err != nil {
			t.Fatal(err)
		}
		out[it.Pos] = sn
	}
	return out
}

func wrongKey(sn snapshot) string {
	for _, o := range sn.Options {
		if o.Key != sn.AnswerKey {
			return o.Key
		}
	}
	return ""
}

// answerK answers the first k positions correctly and the rest wrongly.
func (e *env) answerK(t *testing.T, v ExamView, k int) {
	t.Helper()
	keys := e.keys(t, v.ID)
	for _, it := range v.Items {
		ans := wrongKey(keys[it.Pos])
		if it.Pos <= k {
			ans = keys[it.Pos].AnswerKey
		}
		if _, err := e.svc.Answer(context.Background(), e.user, v.ID, it.Pos, ans, 1000, false, 1); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreatePartialTIU(t *testing.T) {
	e := newEnv(t)
	v := e.create(t)
	if len(v.Items) != 35 || v.Status != store.StatusInProgress || len(v.Sections) != 1 || v.Sections[0] != "TIU" {
		t.Fatalf("view %+v", v)
	}
	// 100 min * 35/110 = 1909 s; basis must say it is prorated.
	if v.DurationSeconds != 1909 || v.DurationBasis != BasisProrated {
		t.Errorf("duration %d basis %q", v.DurationSeconds, v.DurationBasis)
	}
	if got := int(v.DeadlineAt.Sub(v.StartedAt) / time.Second); got != 1909 {
		t.Errorf("deadline - start = %d s", got)
	}
	seen := map[string]bool{}
	for i, it := range v.Items {
		if it.Pos != i+1 || it.Section != "TIU" || len(it.Options) < 2 {
			t.Fatalf("item %d: %+v", i, it)
		}
		if seen[it.Stem] {
			t.Errorf("repeated stem in one exam: %q", it.Stem)
		}
		seen[it.Stem] = true
	}
}

func TestGetHidesKeysAndSnapshotIsStable(t *testing.T) {
	e := newEnv(t)
	v := e.create(t)
	e.answerK(t, v, 3)
	var raws [2]string
	for i := range raws {
		g, err := e.svc.Get(context.Background(), e.user, v.ID)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(g)
		raws[i] = string(b)
	}
	if raws[0] != raws[1] {
		t.Error("two GETs of the same exam differ")
	}
	for _, banned := range []string{"answer_key", "steps", `"correct"`, "misconception"} {
		if strings.Contains(raws[0], banned) {
			t.Errorf("open exam JSON leaks %q", banned)
		}
	}
	if !strings.Contains(raws[0], `"answer"`) {
		t.Error("own answers missing from the open exam view")
	}
}

func TestScoreMatchesHandCountAndScoreSection(t *testing.T) {
	spc, _ := spec.Load()
	sp, _ := spc.Get("skd-2024")
	sec, _ := sp.Section("TIU")
	track, _ := sp.Track("umum")
	for _, k := range []int{0, 1, 15, 16, 35} {
		e := newEnv(t)
		v := e.create(t)
		e.answerK(t, v, k)
		res, err := e.svc.Submit(context.Background(), e.user, v.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Sections) != 1 || !res.Partial || res.Status != store.StatusSubmitted {
			t.Fatalf("k=%d: result %+v", k, res)
		}
		s := res.Sections[0]
		if s.Score != 5*k || s.Answered != 35 || s.Max != 175 || res.Total != 5*k {
			t.Errorf("k=%d: section %+v total %d, want score %d", k, s, res.Total, 5*k)
		}
		// 80 is the TIU threshold for track umum: 16 correct = 80 passes.
		if s.Min == nil || *s.Min != 80 || s.Met != (k >= 16) {
			t.Errorf("k=%d: threshold %+v", k, s)
		}
		rs := make([]scoring.Response, 35)
		for i := range rs {
			rs[i] = scoring.Response{Section: "TIU", Answered: true, Correct: i < k}
		}
		want, err := scoring.ScoreSection(sec, track, rs)
		if err != nil {
			t.Fatal(err)
		}
		if want.Score != s.Score || want.Met != s.Met {
			t.Errorf("k=%d: ScoreSection %+v vs result %+v", k, want, s)
		}
	}
}

func TestSeqHighestWins(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	keys := e.keys(t, v.ID)
	right, wrong := keys[1].AnswerKey, wrongKey(keys[1])

	r, err := e.svc.Answer(ctx, e.user, v.ID, 1, right, 100, true, 2)
	if err != nil || r.Answer != right || r.Seq != 2 || !r.Flagged {
		t.Fatalf("seq 2: %+v %v", r, err)
	}
	// An older write that arrives late must not overwrite, and reports the stored value.
	r, err = e.svc.Answer(ctx, e.user, v.ID, 1, wrong, 200, false, 1)
	if err != nil || r.Answer != right || r.Seq != 2 || r.TimeMS != 100 || !r.Flagged {
		t.Fatalf("seq 1 after 2: %+v %v", r, err)
	}
	// Same seq is a retry: ignored.
	r, _ = e.svc.Answer(ctx, e.user, v.ID, 1, wrong, 200, false, 2)
	if r.Answer != right {
		t.Errorf("equal seq replaced the answer: %+v", r)
	}
	// The answer can be changed with a higher seq.
	r, _ = e.svc.Answer(ctx, e.user, v.ID, 1, wrong, 300, false, 3)
	if r.Answer != wrong || r.Seq != 3 || r.Flagged {
		t.Errorf("seq 3: %+v", r)
	}
	res, err := e.svc.Submit(ctx, e.user, v.ID)
	if err != nil || res.Sections[0].Score != 0 {
		t.Errorf("changed-to-wrong answer scored %+v %v", res, err)
	}
}

func TestDeadlineGrace(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	keys := e.keys(t, v.ID)

	e.clock = v.DeadlineAt.Add(4 * time.Second)
	if _, err := e.svc.Answer(ctx, e.user, v.ID, 1, keys[1].AnswerKey, 1, false, 1); err != nil {
		t.Fatalf("PUT at deadline+4s: %v", err)
	}
	e.clock = v.DeadlineAt.Add(6 * time.Second)
	if _, err := e.svc.Answer(ctx, e.user, v.ID, 2, keys[2].AnswerKey, 1, false, 1); !errors.Is(err, ErrClosed) {
		t.Fatalf("PUT at deadline+6s: %v, want ErrClosed", err)
	}
	if _, err := e.svc.Submit(ctx, e.user, v.ID); err != nil {
		t.Fatal(err)
	}
}

func TestExpiredOnReadScoresBlankAsZero(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	keys := e.keys(t, v.ID)
	for pos := 1; pos <= 3; pos++ {
		if _, err := e.svc.Answer(ctx, e.user, v.ID, pos, keys[pos].AnswerKey, 1, false, 1); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := e.svc.Result(ctx, e.user, v.ID); !errors.Is(err, ErrInProgress) {
		t.Fatalf("result while open: %v", err)
	}

	e.clock = v.DeadlineAt.Add(time.Minute)
	g, err := e.svc.Get(ctx, e.user, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if g.Status != store.StatusExpired || g.SubmittedAt == nil || !g.SubmittedAt.Equal(v.DeadlineAt) {
		t.Fatalf("status %q submitted_at %v", g.Status, g.SubmittedAt)
	}
	if g.Items[0].AnswerKey == "" || len(g.Items[0].Steps) == 0 || g.Items[0].Response == nil || g.Items[0].Response.Correct == nil {
		t.Errorf("closed exam should reveal key, steps and correctness: %+v", g.Items[0])
	}
	res, err := e.svc.Result(ctx, e.user, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != store.StatusExpired || res.Sections[0].Score != 15 || res.Sections[0].Answered != 3 {
		t.Errorf("expired result %+v", res)
	}
}

func TestSubmitTwiceIsSafe(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	e.answerK(t, v, 10)
	a, err := e.svc.Submit(ctx, e.user, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	e.clock = e.clock.Add(time.Second)
	b, err := e.svc.Submit(ctx, e.user, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if a.Sections[0].Score != 50 || a.Total != b.Total || !a.SubmittedAt.Equal(b.SubmittedAt) {
		t.Errorf("first %+v second %+v", a, b)
	}
	keys := e.keys(t, v.ID)
	if _, err := e.svc.Answer(ctx, e.user, v.ID, 1, keys[1].AnswerKey, 1, false, 9); !errors.Is(err, ErrClosed) {
		t.Errorf("answer after submit: %v", err)
	}
}

func TestCreateValidation(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	cases := []struct {
		name, spec, track string
		sections          []string
	}{
		{"TWK has no source", "skd-2024", "umum", []string{"TWK"}},
		{"TKP has no source", "skd-2024", "umum", []string{"TKP"}},
		{"TIU plus TWK", "skd-2024", "umum", []string{"TIU", "TWK"}},
		{"unknown spec", "nope", "umum", []string{"TIU"}},
		{"unknown track", "skd-2024", "nope", []string{"TIU"}},
		{"unknown section", "skd-2024", "umum", []string{"XYZ"}},
		{"no sections", "skd-2024", "umum", nil},
		{"duplicate section", "skd-2024", "umum", []string{"TIU", "TIU"}},
	}
	for _, c := range cases {
		if _, err := e.svc.Create(ctx, e.user, c.spec, c.track, c.sections); !errors.Is(err, ErrInvalid) {
			t.Errorf("%s: err %v, want ErrInvalid", c.name, err)
		}
	}
}

func TestAnswerValidation(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	cases := []struct {
		name string
		pos  int
		ans  string
		ms   int
		seq  int
		want error
	}{
		{"bad letter", 1, "Z", 1, 1, ErrInvalid},
		{"empty", 1, "", 1, 1, ErrInvalid},
		{"negative time", 1, "A", -1, 1, ErrInvalid},
		{"negative seq", 1, "A", 1, -1, ErrInvalid},
		{"pos out of range", 99, "A", 1, 1, ErrNotFound},
	}
	for _, c := range cases {
		if _, err := e.svc.Answer(ctx, e.user, v.ID, c.pos, c.ans, c.ms, false, c.seq); !errors.Is(err, c.want) {
			t.Errorf("%s: err %v, want %v", c.name, err, c.want)
		}
	}
	if _, err := e.svc.Get(ctx, e.user, "not-a-uuid"); !errors.Is(err, ErrNotFound) {
		t.Errorf("bad id: %v", err)
	}
}

func TestExamAndDrillAttemptsAreSeparate(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	if _, err := drill.New(e.st).Get(ctx, v.ID); !errors.Is(err, drill.ErrNotFound) {
		t.Errorf("drill route sees an exam: %v", err)
	}
	d, err := drill.New(e.st).Create(ctx, persentase.PatternID, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.svc.Get(ctx, e.user, d.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("exam route sees a drill attempt: %v", err)
	}
}

func TestSpecsMarkAvailability(t *testing.T) {
	e := newEnv(t)
	got := map[string]bool{}
	for _, sp := range e.svc.Specs() {
		if sp.ID != "skd-2024" {
			continue
		}
		for _, sec := range sp.Sections {
			got[sec.ID] = sec.Available
		}
	}
	if !got["TIU"] || got["TWK"] || got["TKP"] {
		t.Errorf("availability %v", got)
	}
}

func TestExamBelongsToItsUser(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	v := e.create(t)
	other := e.newUser(t)
	keys := e.keys(t, v.ID)

	if _, err := e.svc.Get(ctx, other, v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get by another user: %v, want ErrNotFound", err)
	}
	if _, err := e.svc.Answer(ctx, other, v.ID, 1, keys[1].AnswerKey, 1, false, 1); !errors.Is(err, ErrNotFound) {
		t.Errorf("Answer by another user: %v, want ErrNotFound", err)
	}
	if _, err := e.svc.Submit(ctx, other, v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Submit by another user: %v, want ErrNotFound", err)
	}
	if _, err := e.svc.Result(ctx, other, v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Result by another user: %v, want ErrNotFound", err)
	}
	if _, err := e.svc.Get(ctx, "", v.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get with empty user: %v, want ErrNotFound", err)
	}
	// The other user's attempts did not close or change the owner's exam.
	g, err := e.svc.Get(ctx, e.user, v.ID)
	if err != nil || g.Status != store.StatusInProgress {
		t.Errorf("owner Get: %v status=%q", err, g.Status)
	}
}
