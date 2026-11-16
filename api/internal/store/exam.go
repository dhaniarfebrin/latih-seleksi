package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"latihseleksi/api/internal/question"
)

// CreateExam inserts a simulation attempt and its frozen items in one
// transaction. StartedAt and DeadlineAt come from the caller's clock.
func (s *Store) CreateExam(ctx context.Context, e question.Exam, items []question.Item) (question.Exam, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return question.Exam{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = tx.QueryRow(ctx,
		`insert into attempt (mode, user_id, spec_id, spec_body, track, sections, started_at, deadline_at, status)
		 values ('simulasi', $7, $1, $2, $3, $4, $5, $6, 'in_progress') returning id::text`,
		e.SpecID, []byte(e.SpecBody), e.Track, e.Sections, e.StartedAt, e.DeadlineAt, e.UserID).Scan(&e.ID)
	if err != nil {
		return question.Exam{}, fmt.Errorf("insert exam: %w", err)
	}
	e.Status = "in_progress" // matches the value inserted above
	for _, it := range items {
		if _, err := tx.Exec(ctx,
			`insert into attempt_item (attempt_id, pos, seed, item_snapshot, section, pattern_id) values ($1, $2, $3, $4, $5, $6)`,
			e.ID, it.Pos, it.Seed, []byte(it.Snapshot), it.Section, it.PatternID); err != nil {
			return question.Exam{}, fmt.Errorf("insert item %d: %w", it.Pos, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return question.Exam{}, err
	}
	return e, nil
}

// GetExam returns the simulation attempt with its items and responses, both
// ordered by pos. Drill attempts are not found here.
func (s *Store) GetExam(ctx context.Context, id string) (question.Exam, []question.Item, []question.Response, error) {
	var e question.Exam
	var raw []byte
	err := s.pool.QueryRow(ctx,
		`select id::text, coalesce(user_id::text, ''), spec_id, spec_body, track, sections, started_at, deadline_at, submitted_at, status
		 from attempt where id = $1 and mode = 'simulasi'`, id).
		Scan(&e.ID, &e.UserID, &e.SpecID, &raw, &e.Track, &e.Sections, &e.StartedAt, &e.DeadlineAt, &e.SubmittedAt, &e.Status)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return question.Exam{}, nil, nil, question.ErrNotFound
	}
	if err != nil {
		return question.Exam{}, nil, nil, err
	}
	e.SpecBody = raw

	rows, err := s.pool.Query(ctx,
		`select pos, seed, item_snapshot, coalesce(section, ''), coalesce(pattern_id, '')
		 from attempt_item where attempt_id = $1 order by pos`, id)
	if err != nil {
		return question.Exam{}, nil, nil, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (question.Item, error) {
		var it question.Item
		var snap []byte
		err := r.Scan(&it.Pos, &it.Seed, &snap, &it.Section, &it.PatternID)
		it.Snapshot = snap
		return it, err
	})
	if err != nil {
		return question.Exam{}, nil, nil, err
	}

	rows, err = s.pool.Query(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at, seq, flagged
		 from response where attempt_id = $1 order by pos`, id)
	if err != nil {
		return question.Exam{}, nil, nil, err
	}
	resps, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (question.Response, error) {
		var x question.Response
		err := r.Scan(&x.Pos, &x.Answer, &x.Correct, &x.TimeMS, &x.MisconceptionID, &x.AnsweredAt, &x.Seq, &x.Flagged)
		return x, err
	})
	if err != nil {
		return question.Exam{}, nil, nil, err
	}
	return e, items, resps, nil
}

// ListExamIDs returns the ids of the user's simulation attempts, newest first.
func (s *Store) ListExamIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	rows, err := s.pool.Query(ctx,
		`select id::text from attempt where user_id = $1 and mode = 'simulasi'
		 order by started_at desc, id limit $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

// SaveExamResponse stores r unless a response with an equal or higher seq is
// already stored for the position (highest seq wins). It returns the stored
// response either way.
func (s *Store) SaveExamResponse(ctx context.Context, id string, r question.Response) (question.Response, error) {
	var mis any
	if r.MisconceptionID != "" {
		mis = r.MisconceptionID
	}
	if _, err := s.pool.Exec(ctx,
		`insert into response (attempt_id, pos, answer, correct, time_ms, misconception_id, seq, flagged)
		 values ($1, $2, $3, $4, $5, $6, $7, $8)
		 on conflict (attempt_id, pos) do update set
		   answer = excluded.answer, correct = excluded.correct, time_ms = excluded.time_ms,
		   misconception_id = excluded.misconception_id, seq = excluded.seq,
		   flagged = excluded.flagged, answered_at = now()
		 where excluded.seq > response.seq`,
		id, r.Pos, r.Answer, r.Correct, r.TimeMS, mis, r.Seq, r.Flagged); err != nil {
		return question.Response{}, err
	}
	var out question.Response
	err := s.pool.QueryRow(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at, seq, flagged
		 from response where attempt_id = $1 and pos = $2`, id, r.Pos).
		Scan(&out.Pos, &out.Answer, &out.Correct, &out.TimeMS, &out.MisconceptionID, &out.AnsweredAt, &out.Seq, &out.Flagged)
	return out, err
}

// CloseExam moves an in-progress exam to status (submitted or expired) at
// time at. It does nothing if the exam is already closed.
func (s *Store) CloseExam(ctx context.Context, id, status string, at time.Time) error {
	_, err := s.pool.Exec(ctx,
		`update attempt set status = $2, submitted_at = $3
		 where id = $1 and mode = 'simulasi' and status = 'in_progress'`, id, status, at)
	return err
}

// MedianCorrectMS returns the median time_ms of the user's correct,
// non-weighted answers per pattern id, for the given patterns that have at
// least minCorrect such answers. Pattern id and weight test match ResponseFacts.
func (s *Store) MedianCorrectMS(ctx context.Context, userID string, patternIDs []string, minCorrect int) (map[string]int, error) {
	rows, err := s.pool.Query(ctx,
		`select coalesce(i.pattern_id, a.pattern_id, '') as pid,
		        round(percentile_cont(0.5) within group (order by r.time_ms))::int
		 from response r
		 join attempt a on a.id = r.attempt_id
		 left join attempt_item i on i.attempt_id = r.attempt_id and i.pos = r.pos
		 where a.user_id = $1 and r.correct
		   and not coalesce(exists (select 1 from jsonb_array_elements(i.item_snapshot->'options') o where o ? 'weight'), false)
		   and coalesce(i.pattern_id, a.pattern_id, '') = any($2)
		 group by pid
		 having count(*) >= $3`, userID, patternIDs, minCorrect)
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for rows.Next() {
		var id string
		var ms int
		if err := rows.Scan(&id, &ms); err != nil {
			return nil, err
		}
		out[id] = ms
	}
	return out, rows.Err()
}

// SaveReport stores a question report for an item of the exam, replacing an
// earlier report for the same position. An unknown position (foreign key
// violation) or malformed exam id is ErrNotFound.
func (s *Store) SaveReport(ctx context.Context, id string, pos int, reason, note string) error {
	_, err := s.pool.Exec(ctx,
		`insert into question_report (attempt_id, pos, reason, note) values ($1, $2, $3, $4)
		 on conflict (attempt_id, pos) do update set reason = excluded.reason, note = excluded.note, created_at = now()`,
		id, pos, reason, note)
	var pgErr *pgconn.PgError
	if isInvalidUUID(err) || (errors.As(err, &pgErr) && pgErr.Code == "23503") {
		return question.ErrNotFound
	}
	return err
}
