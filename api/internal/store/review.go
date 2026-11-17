package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"latihseleksi/api/internal/fsrs"
	"latihseleksi/api/internal/review"
)

var _ review.Store = (*Store)(nil)

// Apply writes the review log and the next card for every outcome that has no
// log row yet for the attempt, in one transaction. A repeated call is a no-op:
// the log primary key (attempt_id, kind, key) decides.
func (s *Store) Apply(ctx context.Context, userID, attemptID string, outcomes []review.Outcome, step review.Step) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	for _, o := range outcomes {
		var prev *fsrs.Card
		var c fsrs.Card
		err := tx.QueryRow(ctx,
			`select stability, difficulty, due, last_review, reps, lapses from review_card
			 where user_id = $1 and kind = $2 and key = $3 for update`, userID, string(o.Kind), o.Key).
			Scan(&c.Stability, &c.Difficulty, &c.Due, &c.LastReview, &c.Reps, &c.Lapses)
		switch {
		case err == nil:
			prev = &c
		case !errors.Is(err, pgx.ErrNoRows):
			return err
		}

		next, l := step(prev, o)
		var before, beforeD any
		if prev != nil {
			before, beforeD = l.StabilityBefore, l.DifficultyBefore
		}
		tag, err := tx.Exec(ctx,
			`insert into review_log (user_id, attempt_id, kind, key, grade, reviewed_at, elapsed_days, stability_before, difficulty_before)
			 values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 on conflict (attempt_id, kind, key) do nothing`,
			userID, attemptID, string(o.Kind), o.Key, int(l.Grade), l.ReviewedAt, l.ElapsedDays, before, beforeD)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			continue
		}
		if _, err := tx.Exec(ctx,
			`insert into review_card (user_id, kind, key, stability, difficulty, due, last_review, reps, lapses)
			 values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 on conflict (user_id, kind, key) do update set stability = excluded.stability,
			   difficulty = excluded.difficulty, due = excluded.due, last_review = excluded.last_review,
			   reps = excluded.reps, lapses = excluded.lapses`,
			userID, string(o.Kind), o.Key, next.Stability, next.Difficulty, next.Due, next.LastReview, next.Reps, next.Lapses); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// Due returns the keys of the user's cards of the kind that are due at now,
// most overdue first.
func (s *Store) Due(ctx context.Context, userID string, kind review.Kind, now time.Time, limit int) ([]string, error) {
	rows, err := s.pool.Query(ctx,
		`select key from review_card where user_id = $1 and kind = $2 and due <= $3
		 order by due, key limit $4`, userID, string(kind), now, limit)
	if err != nil {
		return nil, err
	}
	keys, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if keys == nil {
		keys = []string{}
	}
	return keys, err
}
