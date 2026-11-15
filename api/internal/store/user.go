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

// DeleteUser removes the app_user row; the user's attempts, items and
// responses go with it through on delete cascade. An unknown id is not an error.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `delete from app_user where id = $1`, id)
	return err
}
