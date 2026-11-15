package insight

import "context"

// Store is what Service needs from persistence.
type Store interface {
	// ResponseFacts returns the user's most recent answered responses of drill
	// and simulation attempts, newest first, at most limit rows.
	ResponseFacts(ctx context.Context, userID string, limit int) ([]Fact, error)
}
