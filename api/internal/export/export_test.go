package export

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	d   Data
	err error
}

func (f fakeStore) ExportUser(context.Context, string) (Data, error) { return f.d, f.err }

func TestBuildStampsTimeAndNonNilAttempts(t *testing.T) {
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	s := &Service{st: fakeStore{}, now: func() time.Time { return at }}
	d, err := s.Build(context.Background(), "u")
	if err != nil {
		t.Fatal(err)
	}
	if !d.ExportedAt.Equal(at) || d.Attempts == nil {
		t.Errorf("d = %+v", d)
	}
}

func TestBuildPropagatesError(t *testing.T) {
	boom := errors.New("boom")
	s := New(fakeStore{err: boom})
	if _, err := s.Build(context.Background(), "u"); !errors.Is(err, boom) {
		t.Errorf("err = %v", err)
	}
}
