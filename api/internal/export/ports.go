package export

import "context"

// Store is what Service needs from persistence.
type Store interface {
	// ExportUser returns everything the user owns. ExportedAt is left zero.
	// Answer keys (item snapshots) are included for closed attempts only.
	ExportUser(ctx context.Context, userID string) (Data, error)
}
