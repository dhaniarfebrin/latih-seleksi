// Package spec holds versioned exam specifications as data. A spec is plain
// JSON embedded in the binary; adding a new version means adding a new file.
package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Scoring kinds.
const (
	// KindRightWrong awards Correct, Wrong or Blank points per item (TWK, TIU).
	KindRightWrong = "right-wrong"
	// KindWeighted awards the weight of the chosen option, Min..Max, or Blank (TKP).
	KindWeighted = "weighted"
)

// Spec is one exam specification (for example one KepmenPANRB decree).
type Spec struct {
	ID          string    `json:"id"`
	Label       string    `json:"label"`
	Source      string    `json:"source"`
	SourceURL   string    `json:"source_url"`
	EnactedOn   string    `json:"enacted_on"`
	Scope       string    `json:"scope"`
	DurationMin int       `json:"duration_min"`
	Sections    []Section `json:"sections"`
	Tracks      []Track   `json:"tracks"`
	MaxTotal    int       `json:"max_total"`
}

// Section is one subtest (TWK, TIU, TKP).
type Section struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Items   int      `json:"items"`
	Scoring Scoring  `json:"scoring"`
	Aspects []string `json:"aspects,omitempty"`
}

// Scoring describes how one item of a section is scored. Fields not used by
// Kind must be zero.
type Scoring struct {
	Kind    string `json:"kind"`
	Correct int    `json:"correct,omitempty"`
	Wrong   int    `json:"wrong,omitempty"`
	Min     int    `json:"min,omitempty"`
	Max     int    `json:"max,omitempty"`
	Blank   int    `json:"blank,omitempty"`
}

// Track is a pass rule for one applicant path. A track passes when every
// threshold is met and, if TotalMin is set, the cumulative total is met. A
// section without a threshold in the track cannot fail it.
type Track struct {
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Thresholds map[string]int `json:"thresholds,omitempty"`
	TotalMin   int            `json:"total_min,omitempty"`
	// Weights gives each section its percent of the weighted stage score (section
	// id -> percent, summing to 100). Empty when the track has no weighted score.
	Weights map[string]int `json:"weights,omitempty"`
}

// MaxScore is the highest score the section can give.
func (s Section) MaxScore() int {
	if s.Scoring.Kind == KindWeighted {
		return s.Items * s.Scoring.Max
	}
	return s.Items * s.Scoring.Correct
}

// Parse decodes a spec strictly: unknown fields and trailing data are errors.
// It does not validate; call Validate (Registry does).
func Parse(data []byte) (Spec, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var s Spec
	if err := dec.Decode(&s); err != nil {
		return Spec{}, fmt.Errorf("spec: parse: %w", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return Spec{}, fmt.Errorf("spec: parse: unexpected data after spec")
	}
	return s, nil
}

// Section returns the section with the given id.
func (s Spec) Section(id string) (Section, bool) {
	for _, sec := range s.Sections {
		if sec.ID == id {
			return sec, true
		}
	}
	return Section{}, false
}

// Track returns the track with the given id.
func (s Spec) Track(id string) (Track, bool) {
	for _, t := range s.Tracks {
		if t.ID == id {
			return t, true
		}
	}
	return Track{}, false
}
