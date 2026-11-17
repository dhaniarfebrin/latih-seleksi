// Package drill creates, answers and scores drill attempts.
package drill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/insight"
	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/review"
)

const (
	ModeDrill   = "drill"
	DefaultN    = 8
	DefaultBand = 2
	maxN        = 20
	maxTimeMS   = 24 * 60 * 60 * 1000
	seedRetries = 5
	// reviewFetch is how many due fact ids a review drill reads before the
	// bank drops stale ones and the list is cut to n.
	reviewFetch = 200
)

// ReviewPatternID is the reserved pattern id of a drill of the user's due TWK
// facts; its items carry their own TWK.<aspek> pattern.
const ReviewPatternID = "ULANG"

var (
	ErrNotFound  = question.ErrNotFound
	ErrInvalid   = question.ErrInvalid
	ErrSubmitted = errors.New("attempt already submitted")
)

func invalid(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, a...))
}

// Pattern describes an available question family.
type Pattern struct {
	ID             string            `json:"id"`
	Label          string            `json:"label"`
	Misconceptions map[string]string `json:"misconceptions"`
}

// Views returned to clients. Key, steps and misconception appear only inside a
// ResponseView, which exists once the position has been answered.
type OptionView struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type ResponseView struct {
	Answer          string    `json:"answer"`
	Correct         bool      `json:"correct"`
	TimeMS          int       `json:"time_ms"`
	AnswerKey       string    `json:"answer_key"`
	Steps           []string  `json:"steps"`
	MisconceptionID string    `json:"misconception_id,omitempty"`
	AnsweredAt      time.Time `json:"answered_at"`
}

type ItemView struct {
	Pos      int           `json:"pos"`
	Stem     string        `json:"stem"`
	Options  []OptionView  `json:"options"`
	Response *ResponseView `json:"response,omitempty"`
}

type AttemptView struct {
	ID          string     `json:"id"`
	Mode        string     `json:"mode"`
	PatternID   string     `json:"pattern_id"`
	Band        int        `json:"band"`
	N           int        `json:"n"`
	StartedAt   time.Time  `json:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Items       []ItemView `json:"items"`
}

type Result struct {
	AttemptID     string         `json:"attempt_id"`
	N             int            `json:"n"`
	Answered      int            `json:"answered"`
	Score         int            `json:"score"`
	AvgTimeMS     int            `json:"avg_time_ms"`
	Misconception map[string]int `json:"misconceptions"`
	SubmittedAt   time.Time      `json:"submitted_at"`
}

type Service struct {
	store    Store
	ordered  []gen.Family // catalog order, for listing
	families map[string]gen.Family
	facts    FactSource // may be nil
	reviewer Reviewer   // may be nil
}

// New builds the service over st and the registered question families
// (catalog order is the listing order). facts adds banked patterns after them
// and may be nil, as may rev (no review scheduling, no ULANG drill).
func New(st Store, ordered []gen.Family, facts FactSource, rev Reviewer) *Service {
	families := make(map[string]gen.Family, len(ordered))
	for _, f := range ordered {
		families[f.ID()] = f
	}
	return &Service{store: st, ordered: ordered, families: families, facts: facts, reviewer: rev}
}

// Patterns lists the registered families in catalog order, then banked patterns.
func (s *Service) Patterns() []Pattern {
	out := make([]Pattern, 0, len(s.ordered))
	for _, f := range s.ordered {
		out = append(out, Pattern{ID: f.ID(), Label: f.Label(), Misconceptions: f.Misconceptions()})
	}
	if s.facts != nil {
		for _, id := range s.facts.PatternIDs() {
			label, _ := s.facts.PatternLabel(id)
			out = append(out, Pattern{ID: id, Label: label, Misconceptions: s.facts.Misconceptions()})
		}
	}
	return out
}

// PatternIDs lists the ids of Patterns in the same order.
func (s *Service) PatternIDs() []string {
	ps := s.Patterns()
	ids := make([]string, len(ps))
	for i, p := range ps {
		ids[i] = p.ID
	}
	return ids
}

// Create generates n questions with random seeds and freezes them for userID.
func (s *Service) Create(ctx context.Context, userID, patternID string, n, band int) (AttemptView, error) {
	if patternID == ReviewPatternID && s.reviewer != nil && s.facts != nil {
		return s.createReview(ctx, userID, n, band)
	}
	fam, generated := s.families[patternID]
	banked := false
	if !generated && s.facts != nil {
		_, banked = s.facts.PatternLabel(patternID)
	}
	if !generated && !banked {
		return AttemptView{}, invalid("unknown pattern_id %q", patternID)
	}
	if n < 1 || n > maxN {
		return AttemptView{}, invalid("n must be between 1 and %d", maxN)
	}
	if band < 1 || band > 3 {
		return AttemptView{}, invalid("band must be between 1 and 3")
	}

	var items []question.Item
	var err error
	if banked {
		items, err = s.drawFacts(patternID, n)
		n = len(items) // a small pool gives fewer than asked
	} else {
		items, err = s.generate(fam, n, band)
	}
	if err != nil {
		return AttemptView{}, err
	}

	a, err := s.store.CreateAttempt(ctx, question.Attempt{UserID: userID, Mode: ModeDrill, PatternID: patternID, Band: band, N: n}, items)
	if err != nil {
		return AttemptView{}, err
	}
	return s.Get(ctx, userID, a.ID)
}

// createReview freezes up to n of the user's due TWK facts, most overdue first.
// Items keep their real aspect pattern and section so the dashboard stays correct.
func (s *Service) createReview(ctx context.Context, userID string, n, band int) (AttemptView, error) {
	if n < 1 || n > maxN {
		return AttemptView{}, invalid("n must be between 1 and %d", maxN)
	}
	if band < 1 || band > 3 {
		return AttemptView{}, invalid("band must be between 1 and 3")
	}
	// Fetch more than n: ByID drops ids no longer reviewed in the bank, and
	// stale ids among the most overdue must not hide valid ones behind them.
	ids, err := s.reviewer.DueFacts(ctx, userID, reviewFetch)
	if err != nil {
		return AttemptView{}, err
	}
	seed, err := question.RandomSeed()
	if err != nil {
		return AttemptView{}, err
	}
	drafts := s.facts.ByID(ids, seed)
	if len(drafts) > n {
		drafts = drafts[:n]
	}
	if len(drafts) == 0 {
		return AttemptView{}, invalid("nothing is due for review")
	}
	items, err := freeze(drafts, true)
	if err != nil {
		return AttemptView{}, err
	}
	a, err := s.store.CreateAttempt(ctx, question.Attempt{UserID: userID, Mode: ModeDrill, PatternID: ReviewPatternID, Band: band, N: len(items)}, items)
	if err != nil {
		return AttemptView{}, err
	}
	return s.Get(ctx, userID, a.ID)
}

// freeze turns drafts into frozen items; withPattern also sets each item's own
// pattern id and section (the part of the id before the first dot).
func freeze(drafts []gen.Draft, withPattern bool) ([]question.Item, error) {
	items := make([]question.Item, 0, len(drafts))
	for i, d := range drafts {
		raw, err := json.Marshal(question.FromDraft(d))
		if err != nil {
			return nil, err
		}
		it := question.Item{Pos: i + 1, Seed: d.Seed, Snapshot: raw}
		if withPattern {
			it.PatternID = d.PatternID
			it.Section, _, _ = strings.Cut(d.PatternID, ".")
		}
		items = append(items, it)
	}
	return items, nil
}

// drawFacts freezes up to n banked questions of the pattern.
func (s *Service) drawFacts(patternID string, n int) ([]question.Item, error) {
	seed, err := question.RandomSeed()
	if err != nil {
		return nil, err
	}
	drafts, err := s.facts.Draw(patternID, n, seed)
	if err != nil {
		return nil, err
	}
	if len(drafts) == 0 {
		return nil, fmt.Errorf("no questions for %s", patternID)
	}
	return freeze(drafts, false)
}

// generate builds n frozen questions of a family with random seeds.
func (s *Service) generate(fam gen.Family, n, band int) ([]question.Item, error) {
	items := make([]question.Item, 0, n)
	seen := make(map[string]bool, n)
	for pos := 1; pos <= n; pos++ {
		// Retry on a generation error or a stem already in this attempt; if
		// every retry repeats a stem, keep the last repeat rather than fail.
		var draft gen.Draft
		var err error
		got := false
		for try := 0; try < seedRetries; try++ {
			var seed int64
			if seed, err = question.RandomSeed(); err != nil {
				return nil, err
			}
			d, genErr := fam.Generate(seed, gen.Band(band))
			if genErr != nil {
				err = genErr
				continue
			}
			draft, got, err = d, true, nil
			if !seen[d.Stem] {
				break
			}
		}
		if !got {
			return nil, fmt.Errorf("generate item %d: %w", pos, err)
		}
		seen[draft.Stem] = true
		raw, err := json.Marshal(question.FromDraft(draft))
		if err != nil {
			return nil, err
		}
		items = append(items, question.Item{Pos: pos, Seed: draft.Seed, Snapshot: raw})
	}
	return items, nil
}

func (s *Service) load(ctx context.Context, userID, id string) (question.Attempt, map[int]question.Snapshot, []question.Item, []question.Response, error) {
	a, items, resps, err := s.store.GetAttempt(ctx, userID, id)
	if err != nil {
		return a, nil, nil, nil, err
	}
	snaps := make(map[int]question.Snapshot, len(items))
	for _, it := range items {
		var sn question.Snapshot
		if err := json.Unmarshal(it.Snapshot, &sn); err != nil {
			return a, nil, nil, nil, fmt.Errorf("decode snapshot %d: %w", it.Pos, err)
		}
		snaps[it.Pos] = sn
	}
	return a, snaps, items, resps, nil
}

func responseView(r question.Response, sn question.Snapshot) *ResponseView {
	return &ResponseView{
		Answer:          r.Answer,
		Correct:         r.Correct,
		TimeMS:          r.TimeMS,
		AnswerKey:       sn.AnswerKey,
		Steps:           sn.Steps,
		MisconceptionID: r.MisconceptionID,
		AnsweredAt:      r.AnsweredAt,
	}
}

// Get returns the attempt; answer data is included only for answered positions.
func (s *Service) Get(ctx context.Context, userID, id string) (AttemptView, error) {
	a, snaps, items, resps, err := s.load(ctx, userID, id)
	if err != nil {
		return AttemptView{}, err
	}
	byPos := make(map[int]question.Response, len(resps))
	for _, r := range resps {
		byPos[r.Pos] = r
	}
	v := AttemptView{
		ID: a.ID, Mode: a.Mode, PatternID: a.PatternID, Band: a.Band, N: a.N,
		StartedAt: a.StartedAt, SubmittedAt: a.SubmittedAt, Items: make([]ItemView, 0, len(items)),
	}
	for _, it := range items {
		sn := snaps[it.Pos]
		iv := ItemView{Pos: it.Pos, Stem: sn.Stem, Options: make([]OptionView, 0, len(sn.Options))}
		for _, o := range sn.Options {
			iv.Options = append(iv.Options, OptionView{Key: o.Key, Text: o.Text})
		}
		if r, ok := byPos[it.Pos]; ok {
			iv.Response = responseView(r, sn)
		}
		v.Items = append(v.Items, iv)
	}
	return v, nil
}

// Answer records the first answer for pos and returns the revealed result.
// Repeating the call returns the stored result unchanged.
func (s *Service) Answer(ctx context.Context, userID, id string, pos int, answer string, timeMS int) (ResponseView, error) {
	if !question.ValidAnswer(answer) {
		return ResponseView{}, invalid("answer must be one of A-E")
	}
	if timeMS < 0 || timeMS > maxTimeMS {
		return ResponseView{}, invalid("time_ms must be between 0 and %d", maxTimeMS)
	}
	a, snaps, _, resps, err := s.load(ctx, userID, id)
	if err != nil {
		return ResponseView{}, err
	}
	sn, ok := snaps[pos]
	if !ok {
		return ResponseView{}, ErrNotFound
	}
	for _, r := range resps {
		if r.Pos == pos {
			return *responseView(r, sn), nil
		}
	}
	if a.SubmittedAt != nil {
		return ResponseView{}, ErrSubmitted
	}
	chosen, found := sn.Option(answer)
	if !found {
		return ResponseView{}, invalid("answer %s is not an option of this question", answer)
	}
	stored, err := s.store.SaveResponse(ctx, userID, id, question.Response{
		Pos: pos, Answer: answer, Correct: answer == sn.AnswerKey, TimeMS: timeMS, MisconceptionID: chosen.MisconceptionID,
	})
	if err != nil {
		return ResponseView{}, err
	}
	return *responseView(stored, sn), nil
}

// Submit closes the attempt and returns score, mean time and misconception tally.
// Unanswered positions count as incorrect. Repeating the call is safe.
func (s *Service) Submit(ctx context.Context, userID, id string) (Result, error) {
	submittedAt, err := s.store.MarkSubmitted(ctx, userID, id)
	if err != nil {
		return Result{}, err
	}
	a, snaps, _, resps, err := s.load(ctx, userID, id)
	if err != nil {
		return Result{}, err
	}
	res := Result{AttemptID: a.ID, N: a.N, Answered: len(resps), Misconception: map[string]int{}, SubmittedAt: submittedAt}
	total := 0
	for _, r := range resps {
		if r.Correct {
			res.Score++
		}
		if r.MisconceptionID != "" {
			res.Misconception[r.MisconceptionID]++
		}
		total += r.TimeMS
	}
	if len(resps) > 0 {
		res.AvgTimeMS = total / len(resps)
	}
	if s.reviewer != nil {
		if err := s.reviewer.Record(ctx, userID, id, reviewOutcomes(a, snaps, resps, res.Score)); err != nil {
			return Result{}, err
		}
	}
	return res, nil
}

// reviewOutcomes is what a submitted drill teaches the scheduler. Banked
// (TWK) answers review their fact; a generated drill reviews its pattern, as
// Good when accuracy over the answered items reaches insight.QueueAccuracyPct.
// An attempt with no answers schedules nothing.
func reviewOutcomes(a question.Attempt, snaps map[int]question.Snapshot, resps []question.Response, score int) []review.Outcome {
	var out []review.Outcome
	for _, r := range resps {
		if id := snaps[r.Pos].FactID; id != "" {
			out = append(out, review.Outcome{Kind: review.KindFact, Key: id, Correct: r.Correct})
		}
	}
	if len(out) == 0 && len(resps) > 0 && a.PatternID != ReviewPatternID {
		out = append(out, review.Outcome{Kind: review.KindPattern, Key: a.PatternID, Correct: score*100 >= insight.QueueAccuracyPct*len(resps)})
	}
	return out
}

// Report stores a user's report on the question at pos, replacing an earlier
// one. Only answered positions can be reported, since only then the key is
// visible to the user; an unanswered one is ErrInvalid.
func (s *Service) Report(ctx context.Context, userID, id string, pos int, reason, note string) error {
	if !question.ReportReasons[reason] {
		return invalid("reason must be one of kunci_salah, soal_ambigu, typo, lainnya")
	}
	if utf8.RuneCountInString(note) > question.MaxReportNote {
		return invalid("note must be at most %d characters", question.MaxReportNote)
	}
	_, snaps, _, resps, err := s.load(ctx, userID, id)
	if err != nil {
		return err
	}
	if _, ok := snaps[pos]; !ok {
		return ErrNotFound
	}
	for _, r := range resps {
		if r.Pos == pos {
			return s.store.SaveReport(ctx, id, pos, reason, note)
		}
	}
	return invalid("question %d has not been answered", pos)
}
