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
}
