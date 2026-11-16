package exam

import (
	"context"
	"time"

	"latihseleksi/api/internal/question"
)

// Store is what Service needs from persistence. Implementations return
// question.ErrNotFound for an unknown or malformed exam id.
type Store interface {
	// CreateExam inserts the exam (status in_progress) and its frozen items.
	CreateExam(ctx context.Context, e question.Exam, items []question.Item) (question.Exam, error)
	// GetExam returns the exam with items and responses ordered by pos.
	GetExam(ctx context.Context, id string) (question.Exam, []question.Item, []question.Response, error)
	// ListExamIDs returns the user's exam ids, newest first.
	ListExamIDs(ctx context.Context, userID string, limit int) ([]string, error)
	// SaveExamResponse stores r unless an equal or higher seq is stored (highest seq wins)
	// and returns the stored response.
	SaveExamResponse(ctx context.Context, id string, r question.Response) (question.Response, error)
	// CloseExam moves an in-progress exam to status at time at; closed exams are left alone.
	CloseExam(ctx context.Context, id, status string, at time.Time) error
	// MedianCorrectMS returns, per pattern id, the median time_ms of the user's
	// correct non-weighted answers; patterns with fewer than minCorrect correct
	// answers are left out.
	MedianCorrectMS(ctx context.Context, userID string, patternIDs []string, minCorrect int) (map[string]int, error)
	// SaveReport stores the report for an item of the exam, replacing an earlier
	// one for the same position. An unknown position is question.ErrNotFound.
	SaveReport(ctx context.Context, id string, pos int, reason, note string) error
}

// MinCorrectForMedian is how many correct answers a pattern needs before a
// median time is shown (same bar as the dashboard's cause labels).
const MinCorrectForMedian = 3
