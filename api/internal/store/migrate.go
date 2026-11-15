package store

import (
	"context"
	"fmt"
	"io/fs"

	"latihseleksi/api/db"
)

// migrationLockKey serialises concurrent migration runners (pg advisory lock).
const migrationLockKey = 7272001

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
