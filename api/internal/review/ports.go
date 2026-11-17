package review

import (
	"context"
	"time"

	"latihseleksi/api/internal/fsrs"
)

// Log is the before-state of a review, kept for calibration. Before values
// are zero for a card's first review.
type Log struct {
	Grade            fsrs.Grade
	ReviewedAt       time.Time
	ElapsedDays      int
	StabilityBefore  float64
	DifficultyBefore float64
}

// Step turns a card (nil when the item was never reviewed) and an outcome into
// the next card and the log row of that review.
type Step func(prev *fsrs.Card, o Outcome) (fsrs.Card, Log)

// Store is what Service needs from persistence.
type Store interface {
	// Apply handles every outcome in one transaction: if no log row exists for
	// (attemptID, kind, key) it inserts the log and the card produced by step;
	// otherwise it does nothing for that outcome.
	Apply(ctx context.Context, userID, attemptID string, outcomes []Outcome, step Step) error
	// Due returns keys of the user's cards of the kind with due <= now, most
	// overdue first, at most limit.
	Due(ctx context.Context, userID string, kind Kind, now time.Time, limit int) ([]string, error)
}
