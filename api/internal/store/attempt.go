package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"latihseleksi/api/internal/question"
)

// CreateAttempt inserts the attempt and its frozen items in one transaction.
func (s *Store) CreateAttempt(ctx context.Context, a question.Attempt, items []question.Item) (question.Attempt, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return question.Attempt{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = tx.QueryRow(ctx,
		`insert into attempt (user_id, mode, pattern_id, band, n) values ($1, $2, $3, $4, $5)
		 returning id::text, started_at`,
		a.UserID, a.Mode, a.PatternID, a.Band, a.N).Scan(&a.ID, &a.StartedAt)
	if err != nil {
		return question.Attempt{}, fmt.Errorf("insert attempt: %w", err)
	}
	for _, it := range items {
		if _, err := tx.Exec(ctx,
			`insert into attempt_item (attempt_id, pos, seed, item_snapshot, section, pattern_id)
			 values ($1, $2, $3, $4, nullif($5, ''), nullif($6, ''))`,
			a.ID, it.Pos, it.Seed, []byte(it.Snapshot), it.Section, it.PatternID); err != nil {
			return question.Attempt{}, fmt.Errorf("insert item %d: %w", it.Pos, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return question.Attempt{}, err
	}
	return a, nil
}

// GetAttempt returns the user's drill attempt (simulation attempts and other users' attempts are not found) with its items and responses, both ordered by pos.
func (s *Store) GetAttempt(ctx context.Context, userID, id string) (question.Attempt, []question.Item, []question.Response, error) {
	var a question.Attempt
	err := s.pool.QueryRow(ctx,
		`select id::text, mode, pattern_id, band, n, started_at, submitted_at from attempt where id = $1 and user_id = $2 and mode = 'drill'`, id, userID).
		Scan(&a.ID, &a.Mode, &a.PatternID, &a.Band, &a.N, &a.StartedAt, &a.SubmittedAt)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return question.Attempt{}, nil, nil, question.ErrNotFound
	}
	if err != nil {
		return question.Attempt{}, nil, nil, err
	}

	rows, err := s.pool.Query(ctx,
		`select pos, seed, item_snapshot from attempt_item where attempt_id = $1 order by pos`, id)
	if err != nil {
		return question.Attempt{}, nil, nil, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (question.Item, error) {
		var it question.Item
		var raw []byte
		err := r.Scan(&it.Pos, &it.Seed, &raw)
		it.Snapshot = raw
		return it, err
	})
	if err != nil {
		return question.Attempt{}, nil, nil, err
	}

	resps, err := s.responses(ctx, id)
	if err != nil {
		return question.Attempt{}, nil, nil, err
	}
	return a, items, resps, nil
}

func (s *Store) responses(ctx context.Context, id string) ([]question.Response, error) {
	rows, err := s.pool.Query(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at
		 from response where attempt_id = $1 order by pos`, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (question.Response, error) {
		var x question.Response
		err := r.Scan(&x.Pos, &x.Answer, &x.Correct, &x.TimeMS, &x.MisconceptionID, &x.AnsweredAt)
		return x, err
	})
}

// SaveResponse stores the first response for a position of the user's drill
// attempt; later calls keep the stored one. It returns the stored response.
func (s *Store) SaveResponse(ctx context.Context, userID, id string, r question.Response) (question.Response, error) {
	var mis any
	if r.MisconceptionID != "" {
		mis = r.MisconceptionID
	}
	if _, err := s.pool.Exec(ctx,
		`insert into response (attempt_id, pos, answer, correct, time_ms, misconception_id)
		 select a.id, $3, $4, $5, $6, $7 from attempt a where a.id = $1 and a.user_id = $2 and a.mode = 'drill'
		 on conflict (attempt_id, pos) do nothing`,
		id, userID, r.Pos, r.Answer, r.Correct, r.TimeMS, mis); err != nil {
		if isInvalidUUID(err) {
			return question.Response{}, question.ErrNotFound
		}
		return question.Response{}, err
	}
	var out question.Response
	err := s.pool.QueryRow(ctx,
		`select r.pos, r.answer, r.correct, r.time_ms, coalesce(r.misconception_id, ''), r.answered_at
		 from response r join attempt a on a.id = r.attempt_id
		 where r.attempt_id = $1 and r.pos = $2 and a.user_id = $3 and a.mode = 'drill'`, id, r.Pos, userID).
		Scan(&out.Pos, &out.Answer, &out.Correct, &out.TimeMS, &out.MisconceptionID, &out.AnsweredAt)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return question.Response{}, question.ErrNotFound
	}
	return out, err
}

// MarkSubmitted sets submitted_at once on the user's drill attempt and returns the stored value.
func (s *Store) MarkSubmitted(ctx context.Context, userID, id string) (time.Time, error) {
	var t time.Time
	err := s.pool.QueryRow(ctx,
		`update attempt set submitted_at = coalesce(submitted_at, now())
		 where id = $1 and user_id = $2 and mode = 'drill' returning submitted_at`, id, userID).
		Scan(&t)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return time.Time{}, question.ErrNotFound
	}
	return t, err
}
