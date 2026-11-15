package question

import (
	"encoding/json"
	"time"
)

// The records below cross the boundary between the services (drill, exam) and
// the store. They live here so neither side imports the other.

// Attempt is a drill attempt.
type Attempt struct {
	ID          string
	Mode        string
	PatternID   string
	Band        int
	N           int
	StartedAt   time.Time
	SubmittedAt *time.Time
}

// Item is one frozen question of an attempt.
type Item struct {
	Pos      int
	Seed     int64
	Snapshot json.RawMessage
	// Section and PatternID are set for simulation items only.
	Section   string
	PatternID string
}

// Exam is a simulation attempt. SpecBody is the frozen exam spec JSON.
type Exam struct {
	ID          string
	UserID      string
	SpecID      string
	SpecBody    json.RawMessage
	Track       string
	Sections    []string
	StartedAt   time.Time
	DeadlineAt  time.Time
	SubmittedAt *time.Time
	Status      string
}

// Response is a stored answer.
type Response struct {
	Pos             int
	Answer          string
	Correct         bool
	TimeMS          int
	MisconceptionID string
	AnsweredAt      time.Time
	// Seq and Flagged are used by simulation responses only.
	Seq     int
	Flagged bool
}
