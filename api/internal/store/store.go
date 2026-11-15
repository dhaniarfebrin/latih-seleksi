// Package store is the PostgreSQL access layer for drill attempts.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihseleksi/api/db"
)

var ErrNotFound = errors.New("not found")

// migrationLockKey serialises concurrent migration runners (pg advisory lock).
const migrationLockKey = 7272001

type Store struct {
	pool *pgxpool.Pool
}

type Attempt struct {
	ID          string
	Mode        string
	PatternID   string
	Band        int
	N           int
	StartedAt   time.Time
	SubmittedAt *time.Time
}

type Item struct {
	Pos      int
	Seed     int64
	Snapshot json.RawMessage
	// Section and PatternID are set for simulation items only.
	Section   string
	PatternID string
}

// Exam is a simulation attempt. SpecBody is the frozen exam spec JSON.
type Exam struct {
	ID          string
	UserID      string
	SpecID      string
	SpecBody    json.RawMessage
	Track       string
	Sections    []string
	StartedAt   time.Time
	DeadlineAt  time.Time
	SubmittedAt *time.Time
	Status      string
}

type Response struct {
	Pos             int
	Answer          string
	Correct         bool
	TimeMS          int
	MisconceptionID string
	AnsweredAt      time.Time
	// Seq and Flagged are used by simulation responses only.
	Seq     int
	Flagged bool
}

// Open connects to the database. Neon needs sslmode=require in the URL.
func Open(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() { s.pool.Close() }

// Migrate applies embedded migrations not yet recorded in schema_migrations.
func (s *Store) Migrate(ctx context.Context) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `select pg_advisory_lock($1)`, migrationLockKey); err != nil {
		return fmt.Errorf("lock migrations: %w", err)
	}
	defer conn.Exec(context.Background(), `select pg_advisory_unlock($1)`, migrationLockKey) //nolint:errcheck

	if _, err := conn.Exec(ctx, `create table if not exists schema_migrations (
		version text primary key,
		applied_at timestamptz not null default now())`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(db.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	for _, e := range entries {
		version := e.Name()
		var applied bool
		if err := conn.QueryRow(ctx, `select exists (select 1 from schema_migrations where version = $1)`, version).Scan(&applied); err != nil {
			return fmt.Errorf("check %s: %w", version, err)
		}
		if applied {
			continue
		}
		body, err := fs.ReadFile(db.Migrations, "migrations/"+version)
		if err != nil {
			return fmt.Errorf("read %s: %w", version, err)
		}
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin %s: %w", version, err)
		}
		if _, err := tx.Exec(ctx, string(body)); err != nil {
			tx.Rollback(ctx) //nolint:errcheck
			return fmt.Errorf("apply %s: %w", version, err)
		}
		if _, err := tx.Exec(ctx, `insert into schema_migrations (version) values ($1)`, version); err != nil {
			tx.Rollback(ctx) //nolint:errcheck
			return fmt.Errorf("record %s: %w", version, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit %s: %w", version, err)
		}
	}
	return nil
}

// isInvalidUUID reports a malformed uuid parameter (SQLSTATE 22P02).
func isInvalidUUID(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}

// CreateAttempt inserts the attempt and its frozen items in one transaction.
func (s *Store) CreateAttempt(ctx context.Context, a Attempt, items []Item) (Attempt, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Attempt{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = tx.QueryRow(ctx,
		`insert into attempt (mode, pattern_id, band, n) values ($1, $2, $3, $4)
		 returning id::text, started_at`,
		a.Mode, a.PatternID, a.Band, a.N).Scan(&a.ID, &a.StartedAt)
	if err != nil {
		return Attempt{}, fmt.Errorf("insert attempt: %w", err)
	}
	for _, it := range items {
		if _, err := tx.Exec(ctx,
			`insert into attempt_item (attempt_id, pos, seed, item_snapshot) values ($1, $2, $3, $4)`,
			a.ID, it.Pos, it.Seed, []byte(it.Snapshot)); err != nil {
			return Attempt{}, fmt.Errorf("insert item %d: %w", it.Pos, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return Attempt{}, err
	}
	return a, nil
}

// GetAttempt returns the drill attempt (simulation attempts are not found) with its items and responses, both ordered by pos.
func (s *Store) GetAttempt(ctx context.Context, id string) (Attempt, []Item, []Response, error) {
	var a Attempt
	err := s.pool.QueryRow(ctx,
		`select id::text, mode, pattern_id, band, n, started_at, submitted_at from attempt where id = $1 and mode = 'drill'`, id).
		Scan(&a.ID, &a.Mode, &a.PatternID, &a.Band, &a.N, &a.StartedAt, &a.SubmittedAt)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return Attempt{}, nil, nil, ErrNotFound
	}
	if err != nil {
		return Attempt{}, nil, nil, err
	}

	rows, err := s.pool.Query(ctx,
		`select pos, seed, item_snapshot from attempt_item where attempt_id = $1 order by pos`, id)
	if err != nil {
		return Attempt{}, nil, nil, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Item, error) {
		var it Item
		var raw []byte
		err := r.Scan(&it.Pos, &it.Seed, &raw)
		it.Snapshot = raw
		return it, err
	})
	if err != nil {
		return Attempt{}, nil, nil, err
	}

	resps, err := s.responses(ctx, id)
	if err != nil {
		return Attempt{}, nil, nil, err
	}
	return a, items, resps, nil
}

func (s *Store) responses(ctx context.Context, id string) ([]Response, error) {
	rows, err := s.pool.Query(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at
		 from response where attempt_id = $1 order by pos`, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Response, error) {
		var x Response
		err := r.Scan(&x.Pos, &x.Answer, &x.Correct, &x.TimeMS, &x.MisconceptionID, &x.AnsweredAt)
		return x, err
	})
}

// SaveResponse stores the first response for a position; later calls keep the
// stored one. It returns the stored response.
func (s *Store) SaveResponse(ctx context.Context, id string, r Response) (Response, error) {
	var mis any
	if r.MisconceptionID != "" {
		mis = r.MisconceptionID
	}
	if _, err := s.pool.Exec(ctx,
		`insert into response (attempt_id, pos, answer, correct, time_ms, misconception_id)
		 values ($1, $2, $3, $4, $5, $6)
		 on conflict (attempt_id, pos) do nothing`,
		id, r.Pos, r.Answer, r.Correct, r.TimeMS, mis); err != nil {
		return Response{}, err
	}
	var out Response
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
		return time.Time{}, ErrNotFound
	}
	return t, err
}

// UpsertUser returns the id of the app_user with this auth subject (the
// provider's sub), creating it on first sight. Nothing else is stored.
func (s *Store) UpsertUser(ctx context.Context, sub string) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx,
		`insert into app_user (auth_sub) values ($1)
		 on conflict (auth_sub) do update set auth_sub = excluded.auth_sub
		 returning id::text`, sub).Scan(&id)
	return id, err
}

const (
	StatusInProgress = "in_progress"
	StatusSubmitted  = "submitted"
	StatusExpired    = "expired"
)

// CreateExam inserts a simulation attempt and its frozen items in one
// transaction. StartedAt and DeadlineAt come from the caller's clock.
func (s *Store) CreateExam(ctx context.Context, e Exam, items []Item) (Exam, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Exam{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	err = tx.QueryRow(ctx,
		`insert into attempt (mode, user_id, spec_id, spec_body, track, sections, started_at, deadline_at, status)
		 values ('simulasi', $7, $1, $2, $3, $4, $5, $6, 'in_progress') returning id::text`,
		e.SpecID, []byte(e.SpecBody), e.Track, e.Sections, e.StartedAt, e.DeadlineAt, e.UserID).Scan(&e.ID)
	if err != nil {
		return Exam{}, fmt.Errorf("insert exam: %w", err)
	}
	e.Status = StatusInProgress
	for _, it := range items {
		if _, err := tx.Exec(ctx,
			`insert into attempt_item (attempt_id, pos, seed, item_snapshot, section, pattern_id) values ($1, $2, $3, $4, $5, $6)`,
			e.ID, it.Pos, it.Seed, []byte(it.Snapshot), it.Section, it.PatternID); err != nil {
			return Exam{}, fmt.Errorf("insert item %d: %w", it.Pos, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return Exam{}, err
	}
	return e, nil
}

// GetExam returns the simulation attempt with its items and responses, both
// ordered by pos. Drill attempts are not found here.
func (s *Store) GetExam(ctx context.Context, id string) (Exam, []Item, []Response, error) {
	var e Exam
	var raw []byte
	err := s.pool.QueryRow(ctx,
		`select id::text, coalesce(user_id::text, ''), spec_id, spec_body, track, sections, started_at, deadline_at, submitted_at, status
		 from attempt where id = $1 and mode = 'simulasi'`, id).
		Scan(&e.ID, &e.UserID, &e.SpecID, &raw, &e.Track, &e.Sections, &e.StartedAt, &e.DeadlineAt, &e.SubmittedAt, &e.Status)
	if errors.Is(err, pgx.ErrNoRows) || (err != nil && isInvalidUUID(err)) {
		return Exam{}, nil, nil, ErrNotFound
	}
	if err != nil {
		return Exam{}, nil, nil, err
	}
	e.SpecBody = raw

	rows, err := s.pool.Query(ctx,
		`select pos, seed, item_snapshot, coalesce(section, ''), coalesce(pattern_id, '')
		 from attempt_item where attempt_id = $1 order by pos`, id)
	if err != nil {
		return Exam{}, nil, nil, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Item, error) {
		var it Item
		var snap []byte
		err := r.Scan(&it.Pos, &it.Seed, &snap, &it.Section, &it.PatternID)
		it.Snapshot = snap
		return it, err
	})
	if err != nil {
		return Exam{}, nil, nil, err
	}

	rows, err = s.pool.Query(ctx,
		`select pos, answer, correct, time_ms, coalesce(misconception_id, ''), answered_at, seq, flagged
		 from response where attempt_id = $1 order by pos`, id)
	if err != nil {
		return Exam{}, nil, nil, err
	}
	resps, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Response, error) {
		var x Response
		err := r.Scan(&x.Pos, &x.Answer, &x.Correct, &x.TimeMS, &x.MisconceptionID, &x.AnsweredAt, &x.Seq, &x.Flagged)
		return x, err
	})
	if err != nil {
		return Exam{}, nil, nil, err
	}
	return e, items, resps, nil
}

// SaveExamResponse stores r unless a response with an equal or higher seq is
// already stored for the position (highest seq wins). It returns the stored
// response either way.
func (s *Store) SaveExamResponse(ctx context.Context, id string, r Response) (Response, error) {
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
		return Response{}, err
	}
	var out Response
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
