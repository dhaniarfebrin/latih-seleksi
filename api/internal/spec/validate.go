package spec

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Validate checks that the spec is internally consistent. It reports all
// problems found, joined into one error.
func (s Spec) Validate() error {
	var errs []string
	add := func(format string, a ...any) { errs = append(errs, fmt.Sprintf(format, a...)) }

	if s.ID == "" {
		add("id is empty")
	}
	if s.Label == "" {
		add("label is empty")
	}
	if s.Source == "" {
		add("source is empty")
	}
	if s.Scope == "" {
		add("scope is empty")
	}
	if _, err := time.Parse("2006-01-02", s.EnactedOn); err != nil {
		add("enacted_on %q is not a YYYY-MM-DD date", s.EnactedOn)
	}
	if s.DurationMin <= 0 {
		add("duration_min must be > 0, got %d", s.DurationMin)
	}
	if len(s.Sections) == 0 {
		add("no sections")
	}

	sum := 0
	seen := map[string]bool{}
	for _, sec := range s.Sections {
		if sec.ID == "" {
			add("section with empty id")
			continue
		}
		if seen[sec.ID] {
			add("duplicate section id %q", sec.ID)
			continue
		}
		seen[sec.ID] = true
		if sec.Label == "" {
			add("section %s: label is empty", sec.ID)
		}
		if sec.Items <= 0 {
			add("section %s: items must be > 0, got %d", sec.ID, sec.Items)
		}
		validateScoring(sec, add)
		validateAspects(sec, add)
		sum += sec.MaxScore()
	}
	if s.MaxTotal != sum {
		add("max_total is %d, sections add up to %d", s.MaxTotal, sum)
	}

	tracks := map[string]bool{}
	for _, t := range s.Tracks {
		if t.ID == "" {
			add("track with empty id")
			continue
		}
		if tracks[t.ID] {
			add("duplicate track id %q", t.ID)
			continue
		}
		tracks[t.ID] = true
		validateTrack(s, t, add)
	}
	if len(s.Tracks) == 0 {
		add("no tracks")
	}

	if len(errs) > 0 {
		return errors.New("spec " + s.ID + ": " + strings.Join(errs, "; "))
	}
	return nil
}

func validateScoring(sec Section, add func(string, ...any)) {
	sc := sec.Scoring
	switch sc.Kind {
	case KindRightWrong:
		if sc.Correct <= 0 {
			add("section %s: correct must be > 0, got %d", sec.ID, sc.Correct)
		}
		if sc.Wrong < 0 || sc.Blank < 0 {
			add("section %s: wrong and blank must be >= 0", sec.ID)
		}
		if sc.Min != 0 || sc.Max != 0 {
			add("section %s: min/max are not used by %s", sec.ID, KindRightWrong)
		}
	case KindWeighted:
		if sc.Min < 1 || sc.Min > sc.Max {
			add("section %s: weights must satisfy 1 <= min <= max, got %d..%d", sec.ID, sc.Min, sc.Max)
		}
		if sc.Blank < 0 {
			add("section %s: blank must be >= 0", sec.ID)
		}
		if sc.Correct != 0 || sc.Wrong != 0 {
			add("section %s: correct/wrong are not used by %s", sec.ID, KindWeighted)
		}
	default:
		add("section %s: unknown scoring kind %q", sec.ID, sc.Kind)
	}
}

func validateAspects(sec Section, add func(string, ...any)) {
	seen := map[string]bool{}
	for _, a := range sec.Aspects {
		if a == "" {
			add("section %s: empty aspect", sec.ID)
			continue
		}
		if seen[a] {
			add("section %s: duplicate aspect %q", sec.ID, a)
		}
		seen[a] = true
	}
}

func validateTrack(s Spec, t Track, add func(string, ...any)) {
	if t.Label == "" {
		add("track %s: label is empty", t.ID)
	}
	if len(t.Thresholds) == 0 && t.TotalMin == 0 {
		add("track %s: no rule (no thresholds and no total_min)", t.ID)
	}
	ids := make([]string, 0, len(t.Thresholds))
	for id := range t.Thresholds {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		sec, ok := s.Section(id)
		if !ok {
			add("track %s: threshold for unknown section %q", t.ID, id)
			continue
		}
		if th := t.Thresholds[id]; th <= 0 || th > sec.MaxScore() {
			add("track %s: threshold %s=%d must be in 1..%d", t.ID, id, th, sec.MaxScore())
		}
	}
	wids := make([]string, 0, len(t.Weights))
	for id := range t.Weights {
		wids = append(wids, id)
	}
	sort.Strings(wids)
	wsum := 0
	for _, id := range wids {
		if _, ok := s.Section(id); !ok {
			add("track %s: weight for unknown section %q", t.ID, id)
		}
		if w := t.Weights[id]; w <= 0 {
			add("track %s: weight %s=%d must be > 0", t.ID, id, w)
		} else {
			wsum += w
		}
	}
	if len(t.Weights) > 0 && wsum != 100 {
		add("track %s: weights add up to %d, want 100", t.ID, wsum)
	}
	if t.TotalMin < 0 || t.TotalMin > s.MaxTotal {
		add("track %s: total_min %d must be in 1..%d", t.ID, t.TotalMin, s.MaxTotal)
	}
}
