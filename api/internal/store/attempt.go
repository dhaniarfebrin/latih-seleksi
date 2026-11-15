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
		`insert into attempt (mode, pattern_id, band, n) values ($1, $2, $3, $4)
		 returning id::text, started_at`,
		a.Mode, a.PatternID, a.Band, a.N).Scan(&a.ID, &a.StartedAt)
	if err != nil {
		return question.Attempt{}, fmt.Errorf("insert attempt: %w", err)
	}
	for _, it := range items {
		if _, err := tx.Exec(ctx,
			`insert into attempt_item (attempt_id, pos, seed, item_snapshot) values ($1, $2, $3, $4)`,
			a.ID, it.Pos, it.Seed, []byte(it.Snapshot)); err != nil {
			return question.Attempt{}, fmt.Errorf("insert item %d: %w", it.Pos, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return question.Attempt{}, err
	}
	return a, nil
}

// GetAttempt returns the drill attempt (simulation attempts are not found) with its items and responses, both ordered by pos.
func (s *Store) GetAttempt(ctx context.Context, id string) (question.Attempt, []question.Item, []question.Response, error) {
	var a question.Attempt
	err := s.pool.QueryRow(ctx,
		`select id::text, mode, pattern_id, band, n, started_at, submitted_at from attempt where id = $1 and mode = 'drill'`, id).
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

// SaveResponse stores the first response for a position; later calls keep the
// stored one. It returns the stored response.
func (s *Store) SaveResponse(ctx context.Context, id string, r question.Response) (question.Response, error) {
	var mis any
	if r.MisconceptionID != "" {
		mis = r.MisconceptionID
	}
	if _, err := s.pool.Exec(ctx,
		`insert into response (attempt_id, pos, answer, correct, time_ms, misconception_id)
		 values ($1, $2, $3, $4, $5, $6)
		 on conflict (attempt_id, pos) do nothing`,
		id, r.Pos, r.Answer, r.Correct, r.TimeMS, mis); err != nil {
		return question.Response{}, err
	}
	var out question.Response
	err := s.pool.QueryRow(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at
		 from response where attempt_id = $1 and pos = $2`, id, r.Pos).
		Scan(&out.Pos, &out.Answer, &out.Correct, &out.TimeMS, &out.MisconceptionID, &out.AnsweredAt)
	return out, err
}

// MarkSubmitted sets submitted_at once and returns the stored value.
func (s *Store) MarkSubmitted(ctx context.Context, id string) (time.Time, error) {
	var t time.Time
	err := s.pool.QueryRow(ctx,
		`update attempt set submitted_at = coalesce(submitted_at, now()) where id = $1 returning submitted_at`, id).
		Scan(&t)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return time.Time{}, question.ErrNotFound
	}
	return t, err
}
