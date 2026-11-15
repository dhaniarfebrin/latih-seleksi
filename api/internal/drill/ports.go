package drill

import (
	"context"
	"time"

	"latihseleksi/api/internal/question"
)

// Store is what Service needs from persistence. Implementations return
// question.ErrNotFound for an unknown or malformed attempt id.
type Store interface {
	CreateAttempt(ctx context.Context, a question.Attempt, items []question.Item) (question.Attempt, error)
	// GetAttempt returns the drill attempt with items and responses ordered by pos.
	GetAttempt(ctx context.Context, id string) (question.Attempt, []question.Item, []question.Response, error)
	// SaveResponse keeps the first response for a position and returns the stored one.
	SaveResponse(ctx context.Context, id string, r question.Response) (question.Response, error)
	// MarkSubmitted sets submitted_at once and returns the stored value.
	MarkSubmitted(ctx context.Context, id string) (time.Time, error)
}
