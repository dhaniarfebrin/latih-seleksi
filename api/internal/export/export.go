// Package export builds the downloadable copy of one user's data (FR-12).
package export

import (
	"context"
	"encoding/json"
	"time"
)

// Data is the whole export. It never carries the auth subject.
type Data struct {
	ExportedAt time.Time `json:"exported_at"`
	Account    Account   `json:"account"`
	Attempts   []Attempt `json:"attempts"`
	// ReviewCards and ReviewLog are the spaced-repetition state and history.
	ReviewCards []ReviewCard `json:"review_cards"`
	ReviewLog   []ReviewLog  `json:"review_log"`
}

// ReviewCard is the current FSRS state of a TWK fact or TIU pattern.
type ReviewCard struct {
	Kind       string    `json:"kind"`
	Key        string    `json:"key"`
	Stability  float64   `json:"stability"`
	Difficulty float64   `json:"difficulty"`
	Due        time.Time `json:"due"`
	LastReview time.Time `json:"last_review"`
	Reps       int       `json:"reps"`
	Lapses     int       `json:"lapses"`
}

// ReviewLog is one review; the before values are absent for a first review.
type ReviewLog struct {
	AttemptID        string    `json:"attempt_id"`
	Kind             string    `json:"kind"`
	Key              string    `json:"key"`
	Grade            int       `json:"grade"`
	ReviewedAt       time.Time `json:"reviewed_at"`
	ElapsedDays      int       `json:"elapsed_days"`
	StabilityBefore  *float64  `json:"stability_before,omitempty"`
	DifficultyBefore *float64  `json:"difficulty_before,omitempty"`
}

type Account struct {
	CreatedAt time.Time `json:"created_at"`
}

// Attempt is a drill or simulation attempt; fields of the other mode are omitted.
type Attempt struct {
	ID          string     `json:"id"`
	Mode        string     `json:"mode"`
	PatternID   string     `json:"pattern_id,omitempty"`
	Band        *int       `json:"band,omitempty"`
	N           *int       `json:"n,omitempty"`
	SpecID      string     `json:"spec_id,omitempty"`
	Track       string     `json:"track,omitempty"`
	Status      string     `json:"status,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Items       []Item     `json:"items"`
}

// Item is one frozen question. Snapshot holds the answer key, so it is set
// only when the attempt is closed.
type Item struct {
	Pos       int             `json:"pos"`
	Section   string          `json:"section,omitempty"`
	PatternID string          `json:"pattern_id,omitempty"`
	Snapshot  json.RawMessage `json:"snapshot,omitempty"`
	Response  *Response       `json:"response,omitempty"`
	Report    *Report         `json:"report,omitempty"`
}

type Response struct {
	Answer          string    `json:"answer"`
	Correct         bool      `json:"correct"`
	TimeMS          int       `json:"time_ms"`
	MisconceptionID string    `json:"misconception_id,omitempty"`
	Flagged         bool      `json:"flagged"`
	AnsweredAt      time.Time `json:"answered_at"`
}

type Report struct {
	Reason    string    `json:"reason"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Service builds exports.
type Service struct {
	st  Store
	now func() time.Time
}

func New(st Store) *Service { return &Service{st: st, now: time.Now} }

// Build returns the user's data stamped with the export time.
func (s *Service) Build(ctx context.Context, userID string) (Data, error) {
	d, err := s.st.ExportUser(ctx, userID)
	if err != nil {
		return Data{}, err
	}
	if d.Attempts == nil {
		d.Attempts = []Attempt{}
	}
	if d.ReviewCards == nil {
		d.ReviewCards = []ReviewCard{}
	}
	if d.ReviewLog == nil {
		d.ReviewLog = []ReviewLog{}
	}
	d.ExportedAt = s.now().UTC()
	return d, nil
}
