// Package review keeps the user's FSRS cards: a TWK fact or a TIU pattern
// answered in a drill is reviewed once per attempt and rescheduled.
package review

import (
	"context"
	"time"

	"latihseleksi/api/internal/fsrs"
)

// Kind says what a card key refers to.
type Kind string

const (
	// KindFact keys a card by TWK bank entry id.
	KindFact Kind = "fact"
	// KindPattern keys a card by TIU pattern id.
	KindPattern Kind = "pattern"
)

// maxDue bounds the due lists; cards are limited by the bank and the catalog.
const maxDue = 1000

// Outcome is the result of one item in an attempt. Correct maps to Good, wrong
// to Again.
type Outcome struct {
	Kind    Kind
	Key     string
	Correct bool
}

// Counts is what the user has due now.
type Counts struct {
	DueCount    int `json:"due_count"`
	FactsDue    int `json:"facts_due"`
	PatternsDue int `json:"patterns_due"`
}

// Service schedules reviews.
type Service struct {
	st  Store
	now func() time.Time
}

func New(st Store) *Service { return &Service{st: st, now: time.Now} }

// Record reviews each outcome once for the attempt. Calling it again for the
// same attempt changes nothing.
func (s *Service) Record(ctx context.Context, userID, attemptID string, outcomes []Outcome) error {
	if len(outcomes) == 0 {
		return nil
	}
	now := s.now().UTC()
	return s.st.Apply(ctx, userID, attemptID, outcomes, func(prev *fsrs.Card, o Outcome) (fsrs.Card, Log) {
		g := fsrs.Again
		if o.Correct {
			g = fsrs.Good
		}
		var c fsrs.Card
		l := Log{Grade: g, ReviewedAt: now}
		if prev != nil {
			c = *prev
			l.StabilityBefore, l.DifficultyBefore = prev.Stability, prev.Difficulty
			l.ElapsedDays = int(now.Sub(prev.LastReview).Hours() / 24)
			if l.ElapsedDays < 0 {
				l.ElapsedDays = 0
			}
		}
		return fsrs.Review(c, g, now), l
	})
}

// DueFacts lists TWK fact ids due now, most overdue first, at most limit.
func (s *Service) DueFacts(ctx context.Context, userID string, limit int) ([]string, error) {
	return s.st.Due(ctx, userID, KindFact, s.now(), limit)
}

// DuePatterns lists TIU pattern ids due now, most overdue first.
func (s *Service) DuePatterns(ctx context.Context, userID string) ([]string, error) {
	return s.st.Due(ctx, userID, KindPattern, s.now(), maxDue)
}

// Counts returns how many facts and patterns are due now.
func (s *Service) Counts(ctx context.Context, userID string) (Counts, error) {
	facts, err := s.DueFacts(ctx, userID, maxDue)
	if err != nil {
		return Counts{}, err
	}
	pats, err := s.DuePatterns(ctx, userID)
	if err != nil {
		return Counts{}, err
	}
	return Counts{DueCount: len(facts) + len(pats), FactsDue: len(facts), PatternsDue: len(pats)}, nil
}
