package store

import (
	"context"
)

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
