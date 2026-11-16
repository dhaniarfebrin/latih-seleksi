package store

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"latihseleksi/api/internal/question"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	st, err := Open(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return st
}

func TestMedianCorrectMSAndSaveReport(t *testing.T) {
	st := openTest(t)
	ctx := context.Background()
	uid, err := st.UpsertUser(ctx, "store-test-"+strconv.FormatInt(time.Now().UnixNano(), 36))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.DeleteUser(ctx, uid) })

	now := time.Now()
	var items []question.Item
	// pos 1-4 pattern px (correct 100, 200, 400 and one wrong), pos 5-6 pattern py (2 correct).
	for i, p := range []string{"px", "px", "px", "px", "py", "py"} {
		items = append(items, question.Item{Pos: i + 1, Seed: int64(i), Snapshot: []byte(`{"options":[{"key":"A"}]}`), Section: "TIU", PatternID: p})
	}
	e, err := st.CreateExam(ctx, question.Exam{UserID: uid, SpecID: "s", SpecBody: []byte(`{}`), Track: "umum", Sections: []string{"TIU"}, StartedAt: now, DeadlineAt: now.Add(time.Hour)}, items)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range []struct {
		pos     int
		correct bool
		ms      int
	}{{1, true, 100}, {2, true, 400}, {3, true, 200}, {4, false, 9000}, {5, true, 50}, {6, true, 60}} {
		if _, err := st.SaveExamResponse(ctx, e.ID, question.Response{Pos: r.pos, Answer: "A", Correct: r.correct, TimeMS: r.ms, Seq: 1}); err != nil {
			t.Fatal(err)
		}
	}
	got, err := st.MedianCorrectMS(ctx, uid, []string{"px", "py"}, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got["px"] != 200 {
		t.Errorf("median = %v, want only px=200", got)
	}
	if other, err := st.MedianCorrectMS(ctx, "00000000-0000-0000-0000-000000000000", []string{"px"}, 3); err != nil || len(other) != 0 {
		t.Errorf("other user median = %v, %v, want empty", other, err)
	}

	if err := st.SaveReport(ctx, e.ID, 1, "typo", "a"); err != nil {
		t.Fatal(err)
	}
	if err := st.SaveReport(ctx, e.ID, 1, "lainnya", "b"); err != nil {
		t.Fatalf("replace: %v", err)
	}
	var n int
	var reason string
	if err := st.pool.QueryRow(ctx, `select count(*), max(reason) from question_report where attempt_id = $1`, e.ID).Scan(&n, &reason); err != nil || n != 1 || reason != "lainnya" {
		t.Errorf("reports = %d %q, %v, want 1 lainnya", n, reason, err)
	}
	if err := st.SaveReport(ctx, e.ID, 99, "typo", ""); !errors.Is(err, question.ErrNotFound) {
		t.Errorf("unknown pos err = %v, want ErrNotFound", err)
	}
	if err := st.SaveReport(ctx, "not-a-uuid", 1, "typo", ""); !errors.Is(err, question.ErrNotFound) {
		t.Errorf("bad id err = %v, want ErrNotFound", err)
	}
}
