package exam

import (
	"context"
	"fmt"
	"time"

	"latihseleksi/api/internal/question"
)

// fakeStore is an in-memory Store with the same semantics as the PostgreSQL
// one: the highest seq wins, CloseExam only closes an in-progress exam and an
// unknown id is ErrNotFound.
type fakeStore struct {
	exams map[string]*fakeExam
	order []string
}

type fakeExam struct {
	e     question.Exam
	items []question.Item
	resps map[int]question.Response
}

func newFakeStore() *fakeStore { return &fakeStore{exams: map[string]*fakeExam{}} }

var _ Store = (*fakeStore)(nil)

func (f *fakeStore) CreateExam(_ context.Context, e question.Exam, items []question.Item) (question.Exam, error) {
	e.ID = fmt.Sprintf("exam-%d", len(f.order)+1)
	e.Status = StatusInProgress
	f.exams[e.ID] = &fakeExam{e: e, items: items, resps: map[int]question.Response{}}
	f.order = append(f.order, e.ID)
	return e, nil
}

func (f *fakeStore) GetExam(_ context.Context, id string) (question.Exam, []question.Item, []question.Response, error) {
	fe, ok := f.exams[id]
	if !ok {
		return question.Exam{}, nil, nil, question.ErrNotFound
	}
	var resps []question.Response
	for _, it := range fe.items {
		if r, ok := fe.resps[it.Pos]; ok {
			resps = append(resps, r)
		}
	}
	return fe.e, fe.items, resps, nil
}

func (f *fakeStore) ListExamIDs(_ context.Context, userID string, limit int) ([]string, error) {
	var out []string
	for i := len(f.order) - 1; i >= 0 && len(out) < limit; i-- {
		if f.exams[f.order[i]].e.UserID == userID {
			out = append(out, f.order[i])
		}
	}
	return out, nil
}

func (f *fakeStore) SaveExamResponse(_ context.Context, id string, r question.Response) (question.Response, error) {
	fe, ok := f.exams[id]
	if !ok {
		return question.Response{}, question.ErrNotFound
	}
	if old, ok := fe.resps[r.Pos]; ok && r.Seq <= old.Seq {
		return old, nil
	}
	fe.resps[r.Pos] = r
	return r, nil
}

func (f *fakeStore) CloseExam(_ context.Context, id, status string, at time.Time) error {
	fe, ok := f.exams[id]
	if !ok || fe.e.Status != StatusInProgress {
		return nil
	}
	fe.e.Status, fe.e.SubmittedAt = status, &at
	return nil
}
