package exam

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"latihseleksi/api/internal/question"
)

// fakeStore is an in-memory Store with the same semantics as the PostgreSQL
// one: the highest seq wins, CloseExam only closes an in-progress exam and an
// unknown id is ErrNotFound.
type fakeStore struct {
	exams map[string]*fakeExam
	order []string
	// reports by "examID/pos".
	reports map[string]fakeReport
	// getCalls counts GetExam reads, so a test can assert how often a path loads an exam.
	getCalls int
	// getErr, when set, makes GetExam fail like a broken connection.
	getErr error
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
	f.getCalls++
	if f.getErr != nil {
		return question.Exam{}, nil, nil, f.getErr
	}
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

// MedianCorrectMS mirrors the SQL: the user's correct answers per pattern,
// weighted items skipped, upper-middle value for an even count is averaged.
func (f *fakeStore) MedianCorrectMS(_ context.Context, userID string, patternIDs []string, minCorrect int) (map[string]int, error) {
	want := map[string]bool{}
	for _, p := range patternIDs {
		want[p] = true
	}
	times := map[string][]int{}
	for _, id := range f.order {
		fe := f.exams[id]
		if fe.e.UserID != userID {
			continue
		}
		for _, it := range fe.items {
			r, ok := fe.resps[it.Pos]
			if !ok || !r.Correct || !want[it.PatternID] {
				continue
			}
			var sn question.Snapshot
			if err := json.Unmarshal(it.Snapshot, &sn); err != nil {
				return nil, err
			}
			if len(sn.Options) > 0 && sn.Options[0].Weight != 0 {
				continue
			}
			times[it.PatternID] = append(times[it.PatternID], r.TimeMS)
		}
	}
	out := map[string]int{}
	for p, ts := range times {
		if len(ts) < minCorrect {
			continue
		}
		sort.Ints(ts)
		m := len(ts) / 2
		if len(ts)%2 == 1 {
			out[p] = ts[m]
		} else {
			out[p] = (ts[m-1] + ts[m] + 1) / 2
		}
	}
	return out, nil
}

type fakeReport struct{ reason, note string }

func (f *fakeStore) SaveReport(_ context.Context, id string, pos int, reason, note string) error {
	if _, ok := f.exams[id]; !ok {
		return question.ErrNotFound
	}
	if f.reports == nil {
		f.reports = map[string]fakeReport{}
	}
	f.reports[fmt.Sprintf("%s/%d", id, pos)] = fakeReport{reason, note}
	return nil
}
