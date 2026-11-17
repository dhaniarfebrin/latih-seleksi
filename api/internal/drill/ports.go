package drill

import (
	"context"
	"time"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/review"
)

// FactSource serves banked fact questions (TWK) as drill patterns. Pattern ids
// start with the section and a dot, like generated ones ("TWK.integritas").
type FactSource interface {
	PatternIDs() []string
	// PatternLabel returns the user-facing label and false for an unknown id.
	PatternLabel(id string) (string, bool)
	// Misconceptions maps the misconception ids the source's items can emit to labels.
	Misconceptions() map[string]string
	// Draw returns up to n distinct items of the pattern, deterministic per seed.
	Draw(id string, n int, seed int64) ([]gen.Draft, error)
	// ByID returns the items for the given fact ids in that order, each with
	// its drill pattern id in Draft.PatternID; unknown ids are skipped.
	ByID(ids []string, seed int64) []gen.Draft
}

// Reviewer is the spaced-repetition side: it records a submitted drill and
// lists the TWK facts that are due.
type Reviewer interface {
	// Record reviews each outcome once for the attempt; repeating is a no-op.
	Record(ctx context.Context, userID, attemptID string, outcomes []review.Outcome) error
	// DueFacts lists fact ids due now, most overdue first, at most limit.
	DueFacts(ctx context.Context, userID string, limit int) ([]string, error)
}

// Store is what Service needs from persistence. Implementations return
// question.ErrNotFound for an unknown or malformed attempt id and for an
// attempt that belongs to another user.
type Store interface {
	CreateAttempt(ctx context.Context, a question.Attempt, items []question.Item) (question.Attempt, error)
	// GetAttempt returns the user's drill attempt with items and responses ordered by pos.
	GetAttempt(ctx context.Context, userID, id string) (question.Attempt, []question.Item, []question.Response, error)
	// SaveResponse keeps the first response for a position and returns the stored one.
	SaveResponse(ctx context.Context, userID, id string, r question.Response) (question.Response, error)
	// MarkSubmitted sets submitted_at once and returns the stored value.
	MarkSubmitted(ctx context.Context, userID, id string) (time.Time, error)
	// SaveReport stores the report for an item of the attempt, replacing an
	// earlier one. The caller has already checked ownership.
	SaveReport(ctx context.Context, id string, pos int, reason, note string) error
}
