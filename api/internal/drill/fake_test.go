package drill

import (
	"context"
	"errors"
	"fmt"
	"time"

	"latihseleksi/api/internal/question"
)

// fakeStore is an in-memory Store with the same semantics as the PostgreSQL
// one: unknown ids and other users' attempts give ErrNotFound, the first response per position wins and
// submitted_at is set once.
type fakeStore struct {
	attempts map[string]*fakeAttempt
	order    []string
}

type fakeAttempt struct {
	a     question.Attempt
	items []question.Item
	resps map[int]question.Response
}

func newFakeStore() *fakeStore { return &fakeStore{attempts: map[string]*fakeAttempt{}} }

var _ Store = (*fakeStore)(nil)

func (f *fakeStore) CreateAttempt(_ context.Context, a question.Attempt, items []question.Item) (question.Attempt, error) {
	a.ID = fmt.Sprintf("attempt-%d", len(f.order)+1)
	a.StartedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	f.attempts[a.ID] = &fakeAttempt{a: a, items: items, resps: map[int]question.Response{}}
	f.order = append(f.order, a.ID)
	return a, nil
}

func (f *fakeStore) GetAttempt(_ context.Context, userID, id string) (question.Attempt, []question.Item, []question.Response, error) {
	fa, ok := f.attempts[id]
	if !ok || fa.a.UserID != userID {
		return question.Attempt{}, nil, nil, question.ErrNotFound
	}
	var resps []question.Response
	for _, it := range fa.items {
		if r, ok := fa.resps[it.Pos]; ok {
			resps = append(resps, r)
		}
	}
	return fa.a, fa.items, resps, nil
}

func (f *fakeStore) SaveResponse(_ context.Context, userID, id string, r question.Response) (question.Response, error) {
	fa, ok := f.attempts[id]
	if !ok || fa.a.UserID != userID {
		return question.Response{}, question.ErrNotFound
	}
	if old, ok := fa.resps[r.Pos]; ok {
		return old, nil
	}
	r.AnsweredAt = time.Date(2026, 1, 1, 0, 1, 0, 0, time.UTC)
	fa.resps[r.Pos] = r
	return r, nil
}

func (f *fakeStore) MarkSubmitted(_ context.Context, userID, id string) (time.Time, error) {
	fa, ok := f.attempts[id]
	if !ok || fa.a.UserID != userID {
		return time.Time{}, question.ErrNotFound
	}
	if fa.a.SubmittedAt == nil {
		t := time.Date(2026, 1, 1, 0, 2, 0, 0, time.UTC)
		fa.a.SubmittedAt = &t
	}
	return *fa.a.SubmittedAt, nil
}

// errFamilyBroken is what brokenFamily.Generate returns.
var errFamilyBroken = errors.New("family broken")
