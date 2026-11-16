package drill

import (
	"context"
	"time"

	"latihseleksi/api/internal/question"
)

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
