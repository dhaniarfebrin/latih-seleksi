package insight

import "context"

// Store is what Service needs from persistence.
type Store interface {
	// ResponseFacts returns the user's most recent answered responses of drill
	// and simulation attempts, newest first, at most limit rows.
	ResponseFacts(ctx context.Context, userID string, limit int) ([]Fact, error)
}

// PatternSource lists every drill pattern id in display order.
type PatternSource interface {
	PatternIDs() []string
}

// DueSource lists the patterns whose spaced-repetition review is due now.
type DueSource interface {
	DuePatterns(ctx context.Context, userID string) ([]string, error)
}
