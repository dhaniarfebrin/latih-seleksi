// Package insight turns a user's past answers into the weakness dashboard.
package insight

import (
	"context"
	"sort"
	"time"

	"latihseleksi/api/internal/gen"
)

const (
	// FactLimit is how many recent answered responses feed the dashboard.
	FactLimit = 2000
	// MaxSessions caps the session trend list.
	MaxSessions = 10
	// MaxQueue caps the next-drill queue.
	MaxQueue = 3

	// MinCorrectForCause is how many correct answers in a pattern are needed
	// before a wrong answer can be given a cause; fewer gives CauseUnclear.
	MinCorrectForCause = 3
	// SlowFactor: a wrong answer slower than this times the user's median
	// correct time in the pattern is a time-pressure guess.
	SlowFactor = 1.5
	// FastFactor: a wrong answer faster than this times that median is careless.
	FastFactor = 0.5
	// MinAnsweredForQueue is how many answers a pattern needs to enter the queue.
	MinAnsweredForQueue = 3
	// QueueAccuracyPct: patterns below this accuracy are queued.
	QueueAccuracyPct = 80
)

// Cause labels are rough estimates from timing, not diagnoses.
const (
	CauseTimePressure  = "tekanan_waktu"
	CauseCareless      = "ceroboh"
	CauseMisunderstood = "kurang_paham"
	CauseUnclear       = "belum_jelas"
)

// Queue reasons; the UI words them.
const (
	ReasonDue         = "jatuh_tempo"
	ReasonLowAccuracy = "akurasi_rendah"
	ReasonUntried     = "belum_dicoba"
)

// Fact is one answered response with the context needed for statistics.
type Fact struct {
	AttemptID string
	Mode      string
	StartedAt time.Time
	Section   string
	PatternID string
	Correct   bool
	// Weighted marks a weighted (TKP) answer: it has a weight, not right or
	// wrong, so Build ignores it.
	Weighted        bool
	TimeMS          int
	MisconceptionID string
}

// SectionStat is the accuracy and speed of one exam section.
type SectionStat struct {
	Section   string `json:"section"`
	Answered  int    `json:"answered"`
	Correct   int    `json:"correct"`
	AvgTimeMS int    `json:"avg_time_ms"`
}

// CauseCounts is the number of wrong answers per estimated cause.
type CauseCounts struct {
	TimePressure  int `json:"tekanan_waktu"`
	Careless      int `json:"ceroboh"`
	Misunderstood int `json:"kurang_paham"`
	Unclear       int `json:"belum_jelas"`
}

// PatternStat is the accuracy, speed and estimated error causes of one pattern.
type PatternStat struct {
	PatternID string      `json:"pattern_id"`
	Section   string      `json:"section"`
	Answered  int         `json:"answered"`
	Correct   int         `json:"correct"`
	AvgTimeMS int         `json:"avg_time_ms"`
	Causes    CauseCounts `json:"causes"`
	// TopMisconceptionID is the most frequent misconception among wrong answers
	// (ties: smallest id); empty when none was recorded.
	TopMisconceptionID string `json:"top_misconception_id"`
}

// SessionStat is one attempt's result so far.
type SessionStat struct {
	AttemptID string `json:"attempt_id"`
	Mode      string `json:"mode"`
	// PatternID is set for drill attempts (one pattern each); empty for simulations.
	PatternID string    `json:"pattern_id"`
	StartedAt time.Time `json:"started_at"`
	Answered  int       `json:"answered"`
	Correct   int       `json:"correct"`
}

// QueueItem is a suggested next drill. AccuracyPct is set for ReasonLowAccuracy;
// ReasonDue is a pattern whose review is due.
type QueueItem struct {
	PatternID   string `json:"pattern_id"`
	Reason      string `json:"reason"`
	AccuracyPct int    `json:"accuracy_pct"`
}

// Dashboard is the response of GET /v1/me/dashboard. Slices are never nil.
type Dashboard struct {
	Sections []SectionStat `json:"sections"`
	Patterns []PatternStat `json:"patterns"`
	Sessions []SessionStat `json:"sessions"`
	Queue    []QueueItem   `json:"queue"`
}

// Service builds dashboards from stored answers.
type Service struct {
	st       Store
	families []gen.Family
	due      DueSource // may be nil
}

// New returns a Service; families define which patterns can be suggested as
// untried. due lists patterns whose review is due and may be nil.
func New(st Store, families []gen.Family, due DueSource) *Service {
	return &Service{st: st, families: families, due: due}
}

// Dashboard returns the weakness dashboard of the user.
func (s *Service) Dashboard(ctx context.Context, userID string) (Dashboard, error) {
	facts, err := s.st.ResponseFacts(ctx, userID, FactLimit)
	if err != nil {
		return Dashboard{}, err
	}
	ids := make([]string, len(s.families))
	for i, f := range s.families {
		ids[i] = f.ID()
	}
	var due []string
	if s.due != nil {
		if due, err = s.due.DuePatterns(ctx, userID); err != nil {
			return Dashboard{}, err
		}
	}
	return Build(facts, ids, due), nil
}

type patternAcc struct {
	stat    PatternStat
	timeSum int
	correct []int
	wrong   []Fact
	mis     map[string]int
}

// Build is the pure aggregation behind Dashboard. catalog lists all pattern
// ids in display order, used to suggest patterns never tried. due lists the
// patterns whose review is due, most overdue first; they head the queue.
func Build(facts []Fact, catalog, due []string) Dashboard {
	d := Dashboard{
		Sections: []SectionStat{},
		Patterns: []PatternStat{},
		Sessions: []SessionStat{},
		Queue:    []QueueItem{},
	}

	secs := map[string]*SectionStat{}
	secTime := map[string]int{}
	pats := map[string]*patternAcc{}
	sess := map[string]*SessionStat{}
	for _, f := range facts {
		if f.PatternID == "" || f.Weighted {
			continue
		}
		s := secs[f.Section]
		if s == nil {
			s = &SectionStat{Section: f.Section}
			secs[f.Section] = s
		}
		s.Answered++
		secTime[f.Section] += f.TimeMS
		if f.Correct {
			s.Correct++
		}

		p := pats[f.PatternID]
		if p == nil {
			p = &patternAcc{stat: PatternStat{PatternID: f.PatternID, Section: f.Section}, mis: map[string]int{}}
			pats[f.PatternID] = p
		}
		p.stat.Answered++
		p.timeSum += f.TimeMS
		if f.Correct {
			p.stat.Correct++
			p.correct = append(p.correct, f.TimeMS)
		} else {
			p.wrong = append(p.wrong, f)
			if f.MisconceptionID != "" {
				p.mis[f.MisconceptionID]++
			}
		}

		a := sess[f.AttemptID]
		if a == nil {
			a = &SessionStat{AttemptID: f.AttemptID, Mode: f.Mode, StartedAt: f.StartedAt}
			if f.Mode == "drill" {
				a.PatternID = f.PatternID
			}
			sess[f.AttemptID] = a
		}
		a.Answered++
		if f.Correct {
			a.Correct++
		}
	}

	for _, s := range secs {
		s.AvgTimeMS = secTime[s.Section] / s.Answered
		d.Sections = append(d.Sections, *s)
	}
	sort.Slice(d.Sections, func(i, j int) bool { return d.Sections[i].Section < d.Sections[j].Section })

	for _, p := range pats {
		p.stat.AvgTimeMS = p.timeSum / p.stat.Answered
		med := median(p.correct)
		for _, w := range p.wrong {
			switch causeOf(w.TimeMS, med, len(p.correct)) {
			case CauseTimePressure:
				p.stat.Causes.TimePressure++
			case CauseCareless:
				p.stat.Causes.Careless++
			case CauseMisunderstood:
				p.stat.Causes.Misunderstood++
			default:
				p.stat.Causes.Unclear++
			}
		}
		p.stat.TopMisconceptionID = topKey(p.mis)
		d.Patterns = append(d.Patterns, p.stat)
	}
	sort.Slice(d.Patterns, func(i, j int) bool {
		a, b := d.Patterns[i], d.Patterns[j]
		if l, r := a.Correct*b.Answered, b.Correct*a.Answered; l != r {
			return l < r
		}
		return a.PatternID < b.PatternID
	})

	for _, a := range sess {
		d.Sessions = append(d.Sessions, *a)
	}
	sort.Slice(d.Sessions, func(i, j int) bool {
		a, b := d.Sessions[i], d.Sessions[j]
		if !a.StartedAt.Equal(b.StartedAt) {
			return a.StartedAt.After(b.StartedAt)
		}
		return a.AttemptID < b.AttemptID
	})
	if len(d.Sessions) > MaxSessions {
		d.Sessions = d.Sessions[:MaxSessions]
	}

	queued := map[string]bool{}
	for _, id := range due {
		if len(d.Queue) == MaxQueue {
			break
		}
		if !queued[id] {
			queued[id] = true
			d.Queue = append(d.Queue, QueueItem{PatternID: id, Reason: ReasonDue})
		}
	}
	// d.Patterns is already ordered by accuracy, lowest first.
	for _, p := range d.Patterns {
		if len(d.Queue) == MaxQueue {
			break
		}
		if queued[p.PatternID] {
			continue
		}
		if p.Answered >= MinAnsweredForQueue && p.Correct*100 < QueueAccuracyPct*p.Answered {
			queued[p.PatternID] = true
			d.Queue = append(d.Queue, QueueItem{PatternID: p.PatternID, Reason: ReasonLowAccuracy, AccuracyPct: (p.Correct*200 + p.Answered) / (2 * p.Answered)})
		}
	}
	for _, id := range catalog {
		if len(d.Queue) == MaxQueue {
			break
		}
		if pats[id] == nil && !queued[id] {
			d.Queue = append(d.Queue, QueueItem{PatternID: id, Reason: ReasonUntried})
		}
	}
	return d
}

// causeOf estimates why a wrong answer took timeMS, given the median time of
// the user's correct answers in the pattern (nCorrect of them).
func causeOf(timeMS int, med float64, nCorrect int) string {
	if nCorrect < MinCorrectForCause {
		return CauseUnclear
	}
	t := float64(timeMS)
	switch {
	case t > SlowFactor*med:
		return CauseTimePressure
	case t < FastFactor*med:
		return CauseCareless
	default:
		return CauseMisunderstood
	}
}

func median(xs []int) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]int(nil), xs...)
	sort.Ints(s)
	m := len(s) / 2
	if len(s)%2 == 1 {
		return float64(s[m])
	}
	return float64(s[m-1]+s[m]) / 2
}

// topKey returns the key with the highest count, smallest key on ties, or "".
func topKey(m map[string]int) string {
	best, n := "", 0
	for k, c := range m {
		if c > n || (c == n && n > 0 && k < best) {
			best, n = k, c
		}
	}
	return best
}
