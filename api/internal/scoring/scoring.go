// Package scoring computes exam scores from per-item responses against a
// spec. It knows nothing about questions, only about how each response scores.
package scoring

import (
	"fmt"
	"math"

	"latihseleksi/api/internal/spec"
)

// Response is one item's outcome. Correct is used by right-wrong sections,
// Weight (the weight of the chosen option) by weighted sections.
type Response struct {
	Section  string
	Answered bool
	Correct  bool
	Weight   int
}

// SectionResult is the outcome of one subtest. Min is the track threshold, nil
// when the track sets none; Met is true when there is no threshold.
type SectionResult struct {
	ID    string
	Score int
	Max   int
	Min   *int
	Met   bool
}

// Result is the outcome of one attempt against one track.
type Result struct {
	SpecID   string
	Track    string
	Sections []SectionResult
	Total    int
	MaxTotal int
	TotalMin *int
	TotalMet bool
	Passed   bool
}

// Score scores rs against the spec for the given track. Every section must
// have exactly Items responses; unanswered items are Answered=false.
func Score(s spec.Spec, trackID string, rs []Response) (Result, error) {
	track, ok := s.Track(trackID)
	if !ok {
		return Result{}, fmt.Errorf("scoring: unknown track %q in spec %s", trackID, s.ID)
	}

	bySection := map[string][]Response{}
	for i, r := range rs {
		if _, ok := s.Section(r.Section); !ok {
			return Result{}, fmt.Errorf("scoring: response %d: unknown section %q", i, r.Section)
		}
		bySection[r.Section] = append(bySection[r.Section], r)
	}

	res := Result{SpecID: s.ID, Track: track.ID, MaxTotal: s.MaxTotal, TotalMet: true, Passed: true}
	for _, sec := range s.Sections {
		sr, err := ScoreSection(sec, track, bySection[sec.ID])
		if err != nil {
			return Result{}, err
		}
		res.Total += sr.Score
		res.Passed = res.Passed && sr.Met
		res.Sections = append(res.Sections, sr)
	}
	if track.TotalMin > 0 {
		tm := track.TotalMin
		res.TotalMin = &tm
		res.TotalMet = res.Total >= tm
	}
	res.Passed = res.Passed && res.TotalMet
	return res, nil
}

// ScoreSection scores one subtest against the track threshold. It needs
// exactly sec.Items responses, all of this section.
func ScoreSection(sec spec.Section, track spec.Track, rs []Response) (SectionResult, error) {
	if len(rs) != sec.Items {
		return SectionResult{}, fmt.Errorf("scoring: section %s: got %d responses, want %d", sec.ID, len(rs), sec.Items)
	}
	score := 0
	for i, r := range rs {
		if r.Section != sec.ID {
			return SectionResult{}, fmt.Errorf("scoring: section %s: response %d is of section %q", sec.ID, i, r.Section)
		}
		pts, err := itemScore(sec, r)
		if err != nil {
			return SectionResult{}, fmt.Errorf("scoring: response %d: %w", i, err)
		}
		score += pts
	}
	sr := SectionResult{ID: sec.ID, Score: score, Max: sec.MaxScore(), Met: true}
	if th, ok := track.Thresholds[sec.ID]; ok {
		sr.Min = &th
		sr.Met = sr.Score >= th
	}
	return sr, nil
}

func itemScore(sec spec.Section, r Response) (int, error) {
	sc := sec.Scoring
	switch sc.Kind {
	case spec.KindRightWrong:
		switch {
		case !r.Answered:
			return sc.Blank, nil
		case r.Correct:
			return sc.Correct, nil
		default:
			return sc.Wrong, nil
		}
	case spec.KindWeighted:
		if !r.Answered {
			return sc.Blank, nil
		}
		if r.Weight < sc.Min || r.Weight > sc.Max {
			return 0, fmt.Errorf("section %s: weight %d outside %d..%d", sec.ID, r.Weight, sc.Min, sc.Max)
		}
		return r.Weight, nil
	}
	return 0, fmt.Errorf("section %s: unknown scoring kind %q", sec.ID, sc.Kind)
}

// Weighted is the stage score of a track with section weights: the sum of
// weight x score/max over the sections, rounded to 1 decimal (max 100). ok is
// false when the track has no weights or a spec section is missing from rs, so
// a partial attempt has no weighted score.
func Weighted(s spec.Spec, track spec.Track, rs []SectionResult) (score float64, ok bool) {
	if len(track.Weights) == 0 {
		return 0, false
	}
	byID := make(map[string]SectionResult, len(rs))
	for _, r := range rs {
		byID[r.ID] = r
	}
	total := 0.0
	for _, sec := range s.Sections {
		r, found := byID[sec.ID]
		if !found || r.Max <= 0 {
			return 0, false
		}
		total += float64(track.Weights[sec.ID]) * float64(r.Score) / float64(r.Max)
	}
	return math.Round(total*10) / 10, true
}
