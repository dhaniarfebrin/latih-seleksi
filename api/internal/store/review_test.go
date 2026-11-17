package store

import (
	"context"
	"strconv"
	"testing"
	"time"

	"latihseleksi/api/internal/fsrs"
	"latihseleksi/api/internal/review"
)

func TestApplyIsIdempotentPerAttempt(t *testing.T) {
	st := openTest(t)
	ctx := context.Background()
	uid, err := st.UpsertUser(ctx, "review-store-"+strconv.FormatInt(time.Now().UnixNano(), 36))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.DeleteUser(ctx, uid) })

	now := time.Now().UTC().Truncate(time.Microsecond)
	step := func(prev *fsrs.Card, o review.Outcome) (fsrs.Card, review.Log) {
		var c fsrs.Card
		l := review.Log{Grade: fsrs.Good, ReviewedAt: now}
		if prev != nil {
			c = *prev
			l.StabilityBefore, l.DifficultyBefore = prev.Stability, prev.Difficulty
		}
		return fsrs.Review(c, fsrs.Good, now), l
	}
	attempt := "11111111-1111-1111-1111-111111111111"
	out := []review.Outcome{{Kind: review.KindFact, Key: "f1", Correct: true}}
	for i := 0; i < 2; i++ {
		if err := st.Apply(ctx, uid, attempt, out, step); err != nil {
			t.Fatal(err)
		}
	}
	var reps, logs int
	if err := st.pool.QueryRow(ctx, `select reps from review_card where user_id = $1 and key = 'f1'`, uid).Scan(&reps); err != nil || reps != 1 {
		t.Fatalf("reps = %d (%v), want 1", reps, err)
	}
	if err := st.pool.QueryRow(ctx, `select count(*) from review_log where user_id = $1`, uid).Scan(&logs); err != nil || logs != 1 {
		t.Fatalf("logs = %d (%v)", logs, err)
	}
	if err := st.Apply(ctx, uid, "22222222-2222-2222-2222-222222222222", out, step); err != nil {
		t.Fatal(err)
	}
	var before *float64
	if err := st.pool.QueryRow(ctx, `select stability_before from review_log where attempt_id = '22222222-2222-2222-2222-222222222222'`).Scan(&before); err != nil || before == nil {
		t.Errorf("second review lacks stability_before (%v)", err)
	}

	if due, err := st.Due(ctx, uid, review.KindFact, now, 10); err != nil || len(due) != 0 {
		t.Errorf("due now = %v (%v)", due, err)
	}
	if due, err := st.Due(ctx, uid, review.KindFact, now.Add(400*24*time.Hour), 10); err != nil || len(due) != 1 || due[0] != "f1" {
		t.Errorf("due later = %v (%v)", due, err)
	}
	if due, _ := st.Due(ctx, uid, review.KindPattern, now.Add(400*24*time.Hour), 10); len(due) != 0 {
		t.Errorf("pattern due = %v", due)
	}
}
