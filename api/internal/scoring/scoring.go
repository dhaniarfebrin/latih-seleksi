// Package scoring computes exam scores from per-item responses against a
// spec. It knows nothing about questions, only about how each response scores.
package scoring

import (
	"fmt"

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

	scores := map[string]int{}
	counts := map[string]int{}
	for i, r := range rs {
		sec, ok := s.Section(r.Section)
		if !ok {
			return Result{}, fmt.Errorf("scoring: response %d: unknown section %q", i, r.Section)
		}
		pts, err := itemScore(sec, r)
		if err != nil {
			return Result{}, fmt.Errorf("scoring: response %d: %w", i, err)
		}
		scores[sec.ID] += pts
		counts[sec.ID]++
	}

	res := Result{SpecID: s.ID, Track: track.ID, MaxTotal: s.MaxTotal, TotalMet: true, Passed: true}
	for _, sec := range s.Sections {
		if counts[sec.ID] != sec.Items {
			return Result{}, fmt.Errorf("scoring: section %s: got %d responses, want %d", sec.ID, counts[sec.ID], sec.Items)
		}
		sr := SectionResult{ID: sec.ID, Score: scores[sec.ID], Max: sec.MaxScore(), Met: true}
		if th, ok := track.Thresholds[sec.ID]; ok {
			sr.Min = &th
			sr.Met = sr.Score >= th
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
