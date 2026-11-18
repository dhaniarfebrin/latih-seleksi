// Package exam runs simulation attempts: items frozen at creation, a server
// timer, idempotent answers (highest seq wins) and scoring against the
// frozen exam spec. Only sections with an ItemSource can be taken.
package exam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"latihseleksi/api/internal/question"
	"latihseleksi/api/internal/scoring"
	"latihseleksi/api/internal/spec"
)

const (
	// grace is how long after the deadline answers and submits are still
	// accepted, for network delay. Assumption (docs/architecture.md section 8),
	// not an official rule.
	grace     = 5 * time.Second
	maxTimeMS = 24 * 60 * 60 * 1000

	// Duration bases reported to the client.
	BasisSpec     = "spec"
	BasisProrated = "prorated"

	// Exam statuses, as stored.
	StatusInProgress = "in_progress"
	StatusSubmitted  = "submitted"
	StatusExpired    = "expired"
)

var (
	ErrNotFound   = question.ErrNotFound
	ErrInvalid    = question.ErrInvalid
	ErrClosed     = errors.New("exam is closed")
	ErrInProgress = errors.New("exam is still in progress")

	// errCorrupt marks a stored exam whose frozen data cannot be decoded or
	// scored. It is a data problem of one exam, unlike a store error, so List
	// can skip that exam instead of failing the whole list.
	errCorrupt = errors.New("corrupt exam data")
)

func invalid(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, a...))
}

// Views returned to clients. Key, steps and correctness appear only once the
// exam is closed (submitted or expired).
type OptionView struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type ResponseView struct {
	Answer  string `json:"answer"`
	TimeMS  int    `json:"time_ms"`
	Flagged bool   `json:"flagged"`
	Seq     int    `json:"seq"`
	Correct *bool  `json:"correct,omitempty"`
	// Weight is the weight of the chosen option; set only for a closed exam and
	// a weighted section (TKP).
	Weight          *int   `json:"weight,omitempty"`
	MisconceptionID string `json:"misconception_id,omitempty"`
}

type ItemView struct {
	Pos       int           `json:"pos"`
	Section   string        `json:"section"`
	Stem      string        `json:"stem"`
	Options   []OptionView  `json:"options"`
	Response  *ResponseView `json:"response,omitempty"`
	AnswerKey string        `json:"answer_key,omitempty"`
	Steps     []string      `json:"steps,omitempty"`
	// MedianMS is the median time of the user's correct answers in the same
	// pattern; set only for a closed exam, an answered item and a pattern with
	// at least MinCorrectForMedian correct answers.
	MedianMS *int `json:"median_ms,omitempty"`
}

type ExamView struct {
	ID              string     `json:"id"`
	SpecID          string     `json:"spec_id"`
	Track           string     `json:"track"`
	Sections        []string   `json:"sections"`
	Status          string     `json:"status"`
	StartedAt       time.Time  `json:"started_at"`
	DeadlineAt      time.Time  `json:"deadline_at"`
	DurationSeconds int        `json:"duration_seconds"`
	DurationBasis   string     `json:"duration_basis"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	Items           []ItemView `json:"items"`
}

type SectionResult struct {
	ID       string `json:"id"`
	Items    int    `json:"items"`
	Answered int    `json:"answered"`
	Score    int    `json:"score"`
	Max      int    `json:"max"`
	Min      *int   `json:"min,omitempty"`
	Met      bool   `json:"met"`
}

// Result is the outcome per section. Partial is true when not every spec
// section was taken; then there is no overall pass/fail, only each section
// against its threshold.
type Result struct {
	ExamID        string          `json:"exam_id"`
	SpecID        string          `json:"spec_id"`
	Track         string          `json:"track"`
	Status        string          `json:"status"`
	Partial       bool            `json:"partial"`
	DurationBasis string          `json:"duration_basis"`
	Sections      []SectionResult `json:"sections"`
	Total         int             `json:"total"`
	MaxTotal      int             `json:"max_total"`
	SubmittedAt   time.Time       `json:"submitted_at"`
}

type SpecSectionView struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Items     int    `json:"items"`
	Available bool   `json:"available"`
}

type SpecView struct {
	ID          string            `json:"id"`
	Label       string            `json:"label"`
	Scope       string            `json:"scope"`
	DurationMin int               `json:"duration_min"`
	Sections    []SpecSectionView `json:"sections"`
	Tracks      []spec.Track      `json:"tracks"`
}

type Service struct {
	store   Store
	specs   *spec.Registry
	sources map[string]ItemSource // by section id
	now     func() time.Time
}

// New builds the service. A section is available when one of the sources
// serves it; a later source wins for a section two sources serve.
func New(st Store, specs *spec.Registry, sources ...ItemSource) *Service {
	by := map[string]ItemSource{}
	for _, src := range sources {
		for _, sec := range src.Sections() {
			by[sec] = src
		}
	}
	return &Service{store: st, specs: specs, sources: by, now: time.Now}
}

// available reports whether a section can be taken: it has a source.
func (s *Service) available(sec spec.Section) bool {
	_, ok := s.sources[sec.ID]
	return ok
}

// Specs lists the exam specs with which sections can be taken now.
func (s *Service) Specs() []SpecView {
	var out []SpecView
	for _, sp := range s.specs.All() {
		v := SpecView{ID: sp.ID, Label: sp.Label, Scope: sp.Scope, DurationMin: sp.DurationMin, Tracks: sp.Tracks}
		for _, sec := range sp.Sections {
			v.Sections = append(v.Sections, SpecSectionView{ID: sec.ID, Label: sec.Label, Items: sec.Items, Available: s.available(sec)})
		}
		out = append(out, v)
	}
	return out
}

// duration is the time allowed for the chosen sections: the spec duration
// when all sections are taken, else prorated by item count (an assumption,
// not an official figure).
func duration(sp spec.Spec, sections []spec.Section) (time.Duration, string) {
	total, picked := 0, 0
	for _, sec := range sp.Sections {
		total += sec.Items
	}
	for _, sec := range sections {
		picked += sec.Items
	}
	full := time.Duration(sp.DurationMin) * time.Minute
	if picked == total {
		return full, BasisSpec
	}
	return full * time.Duration(picked) / time.Duration(total), BasisProrated
}

// Create picks and freezes the items of the chosen sections and starts the
// timer. Items are numbered in spec section order.
func (s *Service) Create(ctx context.Context, userID, specID, trackID string, sectionIDs []string) (ExamView, error) {
	sp, ok := s.specs.Get(specID)
	if !ok {
		return ExamView{}, invalid("unknown spec_id %q", specID)
	}
	if _, ok := sp.Track(trackID); !ok {
		return ExamView{}, invalid("unknown track %q for spec %s", trackID, sp.ID)
	}
	if len(sectionIDs) == 0 {
		return ExamView{}, invalid("sections must not be empty")
	}
	want := map[string]bool{}
	for _, id := range sectionIDs {
		sec, ok := sp.Section(id)
		if !ok {
			return ExamView{}, invalid("unknown section %q for spec %s", id, sp.ID)
		}
		if want[id] {
			return ExamView{}, invalid("section %q listed twice", id)
		}
		if !s.available(sec) {
			return ExamView{}, invalid("section %q has no question source yet", id)
		}
		want[id] = true
	}

	var chosen []spec.Section
	var items []question.Item
	for _, sec := range sp.Sections {
		if !want[sec.ID] {
			continue
		}
		chosen = append(chosen, sec)
		seed, err := question.RandomSeed()
		if err != nil {
			return ExamView{}, err
		}
		drafts, err := s.sources[sec.ID].Pick(sec.ID, sec.Items, seed)
		if err != nil {
			// A source that cannot serve its section is a server fault, not bad input.
			return ExamView{}, fmt.Errorf("pick items for section %s: %w", sec.ID, err)
		}
		for _, d := range drafts {
			raw, err := json.Marshal(question.FromDraft(d))
			if err != nil {
				return ExamView{}, err
			}
			items = append(items, question.Item{Pos: len(items) + 1, Seed: d.Seed, Snapshot: raw, Section: sec.ID, PatternID: d.PatternID})
		}
	}

	body, err := json.Marshal(sp)
	if err != nil {
		return ExamView{}, err
	}
	dur, _ := duration(sp, chosen)
	now := s.now()
	ids := make([]string, 0, len(chosen))
	for _, sec := range chosen {
		ids = append(ids, sec.ID)
	}
	e, err := s.store.CreateExam(ctx, question.Exam{
		UserID: userID, SpecID: sp.ID, SpecBody: body, Track: trackID, Sections: ids, StartedAt: now, DeadlineAt: now.Add(dur),
	}, items)
	if err != nil {
		return ExamView{}, err
	}
	return s.Get(ctx, userID, e.ID)
}

type loaded struct {
	exam  question.Exam
	spec  spec.Spec
	snaps map[int]question.Snapshot
	items []question.Item
	resps []question.Response
}

func (l loaded) closed() bool { return l.exam.Status != StatusInProgress }

// load reads the exam owned by userID and, if the deadline plus grace has
// passed, closes it as expired first (lazy auto-submit: there is no background
// job). An exam of another user is ErrNotFound, the same as a missing one, so
// ids cannot be probed.
func (s *Service) load(ctx context.Context, userID, id string) (loaded, error) {
	e, items, resps, err := s.store.GetExam(ctx, id)
	if err != nil {
		return loaded{}, err
	}
	if userID == "" || e.UserID != userID {
		return loaded{}, ErrNotFound
	}
	if e.Status == StatusInProgress && s.now().After(e.DeadlineAt.Add(grace)) {
		if err := s.store.CloseExam(ctx, id, StatusExpired, e.DeadlineAt); err != nil {
			return loaded{}, err
		}
		if e, items, resps, err = s.store.GetExam(ctx, id); err != nil {
			return loaded{}, err
		}
	}
	sp, err := spec.Parse(e.SpecBody)
	if err != nil {
		return loaded{}, fmt.Errorf("%w: decode frozen spec: %w", errCorrupt, err)
	}
	snaps := make(map[int]question.Snapshot, len(items))
	for _, it := range items {
		var sn question.Snapshot
		if err := json.Unmarshal(it.Snapshot, &sn); err != nil {
			return loaded{}, fmt.Errorf("%w: decode snapshot %d: %w", errCorrupt, it.Pos, err)
		}
		snaps[it.Pos] = sn
	}
	return loaded{exam: e, spec: sp, snaps: snaps, items: items, resps: resps}, nil
}

func (l loaded) chosenSections() []spec.Section {
	var out []spec.Section
	for _, id := range l.exam.Sections {
		if sec, ok := l.spec.Section(id); ok {
			out = append(out, sec)
		}
	}
	return out
}

// Get returns the exam. Keys and steps are included only once it is closed.
func (s *Service) Get(ctx context.Context, userID, id string) (ExamView, error) {
	l, err := s.load(ctx, userID, id)
	if err != nil {
		return ExamView{}, err
	}
	dur, basis := duration(l.spec, l.chosenSections())
	v := ExamView{
		ID: l.exam.ID, SpecID: l.exam.SpecID, Track: l.exam.Track, Sections: l.exam.Sections,
		Status: l.exam.Status, StartedAt: l.exam.StartedAt, DeadlineAt: l.exam.DeadlineAt,
		DurationSeconds: int(dur / time.Second), DurationBasis: basis, SubmittedAt: l.exam.SubmittedAt,
		Items: make([]ItemView, 0, len(l.items)),
	}
	byPos := make(map[int]question.Response, len(l.resps))
	for _, r := range l.resps {
		byPos[r.Pos] = r
	}
	var medians map[string]int
	if l.closed() {
		var pids []string
		for _, it := range l.items {
			if _, ok := byPos[it.Pos]; ok && it.PatternID != "" {
				pids = append(pids, it.PatternID)
			}
		}
		if len(pids) > 0 {
			if medians, err = s.store.MedianCorrectMS(ctx, userID, pids, MinCorrectForMedian); err != nil {
				return ExamView{}, err
			}
		}
	}
	for _, it := range l.items {
		sn := l.snaps[it.Pos]
		iv := ItemView{Pos: it.Pos, Section: it.Section, Stem: sn.Stem, Options: make([]OptionView, 0, len(sn.Options))}
		for _, o := range sn.Options {
			iv.Options = append(iv.Options, OptionView{Key: o.Key, Text: o.Text})
		}
		if r, ok := byPos[it.Pos]; ok {
			iv.Response = &ResponseView{Answer: r.Answer, TimeMS: r.TimeMS, Flagged: r.Flagged, Seq: r.Seq}
			if l.closed() {
				c := r.Correct
				iv.Response.Correct = &c
				iv.Response.MisconceptionID = r.MisconceptionID
				if sec, ok := l.spec.Section(it.Section); ok && sec.Scoring.Kind == spec.KindWeighted {
					if o, found := sn.Option(r.Answer); found {
						w := o.Weight
						iv.Response.Weight = &w
					}
				}
			}
		}
		if l.closed() {
			iv.AnswerKey, iv.Steps = sn.AnswerKey, sn.Steps
			if _, answered := byPos[it.Pos]; answered {
				if m, ok := medians[it.PatternID]; ok {
					iv.MedianMS = &m
				}
			}
		}
		v.Items = append(v.Items, iv)
	}
	return v, nil
}

// ExamSummary is an exam without items, for lists. Result is set once the exam
// is closed; it carries no answer keys.
type ExamSummary struct {
	ID              string     `json:"id"`
	SpecID          string     `json:"spec_id"`
	Track           string     `json:"track"`
	Sections        []string   `json:"sections"`
	Status          string     `json:"status"`
	StartedAt       time.Time  `json:"started_at"`
	DeadlineAt      time.Time  `json:"deadline_at"`
	DurationSeconds int        `json:"duration_seconds"`
	DurationBasis   string     `json:"duration_basis"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	Result          *Result    `json:"result,omitempty"`
}

const (
	DefaultListLimit = 20
	MaxListLimit     = 50
)

// List returns the user's exams, newest first. Each is loaded through load, so
// an overdue exam is closed (expired) before it is reported. An exam whose
// stored data is corrupt is logged and left out; store errors still fail the call.
func (s *Service) List(ctx context.Context, userID string, limit int) ([]ExamSummary, error) {
	if limit < 1 || limit > MaxListLimit {
		return nil, invalid("limit must be between 1 and %d", MaxListLimit)
	}
	ids, err := s.store.ListExamIDs(ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]ExamSummary, 0, len(ids))
	for _, id := range ids {
		l, err := s.load(ctx, userID, id)
		if errors.Is(err, errCorrupt) {
			slog.Error("exam skipped in list", "exam", id, "err", err)
			continue
		}
		if err != nil {
			return nil, err
		}
		dur, basis := duration(l.spec, l.chosenSections())
		sum := ExamSummary{
			ID: l.exam.ID, SpecID: l.exam.SpecID, Track: l.exam.Track, Sections: l.exam.Sections,
			Status: l.exam.Status, StartedAt: l.exam.StartedAt, DeadlineAt: l.exam.DeadlineAt,
			DurationSeconds: int(dur / time.Second), DurationBasis: basis, SubmittedAt: l.exam.SubmittedAt,
		}
		if l.closed() {
			res, err := l.result()
			if errors.Is(err, errCorrupt) {
				slog.Error("exam skipped in list", "exam", id, "err", err)
				continue
			}
			if err != nil {
				return nil, err
			}
			sum.Result = &res
		}
		out = append(out, sum)
	}
	return out, nil
}

// AnswerInput is one answer submitted for a position.
type AnswerInput struct {
	Pos     int
	Answer  string
	TimeMS  int
	Flagged bool
	Seq     int
}

// Answer stores the answer for in.Pos unless a response with an equal or higher
// seq is already stored; either way it returns the stored response. The answer
// may be changed while the exam is open. Accepted until deadline + grace.
func (s *Service) Answer(ctx context.Context, userID, id string, in AnswerInput) (ResponseView, error) {
	if !question.ValidAnswer(in.Answer) {
		return ResponseView{}, invalid("answer must be one of A-E")
	}
	if in.TimeMS < 0 || in.TimeMS > maxTimeMS {
		return ResponseView{}, invalid("time_ms must be between 0 and %d", maxTimeMS)
	}
	if in.Seq < 0 {
		return ResponseView{}, invalid("seq must not be negative")
	}
	l, err := s.load(ctx, userID, id)
	if err != nil {
		return ResponseView{}, err
	}
	sn, ok := l.snaps[in.Pos]
	if !ok {
		return ResponseView{}, ErrNotFound
	}
	if l.closed() {
		return ResponseView{}, ErrClosed
	}
	chosen, found := sn.Option(in.Answer)
	if !found {
		return ResponseView{}, invalid("answer %s is not an option of this question", in.Answer)
	}
	// Correct means "chose the key option". For a weighted section (TKP) that is
	// only "chose the top weight"; Result scores it by weight, not by Correct.
	stored, err := s.store.SaveExamResponse(ctx, id, question.Response{
		Pos: in.Pos, Answer: in.Answer, Correct: in.Answer == sn.AnswerKey, TimeMS: in.TimeMS,
		MisconceptionID: chosen.MisconceptionID, Seq: in.Seq, Flagged: in.Flagged,
	})
	if err != nil {
		return ResponseView{}, err
	}
	return ResponseView{Answer: stored.Answer, TimeMS: stored.TimeMS, Flagged: stored.Flagged, Seq: stored.Seq}, nil
}

// Report stores a user's report on the question at pos of a closed exam,
// replacing an earlier one for that position. The key and steps are visible
// only after closing, so an exam in progress gives ErrInProgress.
func (s *Service) Report(ctx context.Context, userID, id string, pos int, reason, note string) error {
	if !question.ReportReasons[reason] {
		return invalid("reason must be one of kunci_salah, soal_ambigu, typo, lainnya")
	}
	if utf8.RuneCountInString(note) > question.MaxReportNote {
		return invalid("note must be at most %d characters", question.MaxReportNote)
	}
	l, err := s.load(ctx, userID, id)
	if err != nil {
		return err
	}
	if _, ok := l.snaps[pos]; !ok {
		return ErrNotFound
	}
	if !l.closed() {
		return ErrInProgress
	}
	return s.store.SaveReport(ctx, id, pos, reason, note)
}

// Submit closes the exam and returns the result. Repeating it, or submitting
// an expired exam, returns the same result.
func (s *Service) Submit(ctx context.Context, userID, id string) (Result, error) {
	l, err := s.load(ctx, userID, id)
	if err != nil {
		return Result{}, err
	}
	if !l.closed() {
		if err := s.store.CloseExam(ctx, id, StatusSubmitted, s.now()); err != nil {
			return Result{}, err
		}
	}
	return s.Result(ctx, userID, id)
}

// Result scores a closed exam. An exam still in progress gives ErrInProgress.
func (s *Service) Result(ctx context.Context, userID, id string) (Result, error) {
	l, err := s.load(ctx, userID, id)
	if err != nil {
		return Result{}, err
	}
	return l.result()
}

// result scores an already loaded exam, so callers that hold one (List) do not
// read it a second time. An exam still in progress gives ErrInProgress.
func (l loaded) result() (Result, error) {
	if !l.closed() {
		return Result{}, ErrInProgress
	}
	track, ok := l.spec.Track(l.exam.Track)
	if !ok {
		return Result{}, fmt.Errorf("%w: frozen spec %s lost track %q", errCorrupt, l.spec.ID, l.exam.Track)
	}
	byPos := make(map[int]question.Response, len(l.resps))
	for _, r := range l.resps {
		byPos[r.Pos] = r
	}
	chosen := l.chosenSections()
	_, basis := duration(l.spec, chosen)
	res := Result{
		ExamID: l.exam.ID, SpecID: l.spec.ID, Track: track.ID, Status: l.exam.Status,
		Partial: len(chosen) < len(l.spec.Sections), DurationBasis: basis,
	}
	if l.exam.SubmittedAt != nil {
		res.SubmittedAt = *l.exam.SubmittedAt
	}
	for _, sec := range chosen {
		var rs []scoring.Response
		answered := 0
		for _, it := range l.items {
			if it.Section != sec.ID {
				continue
			}
			r, ok := byPos[it.Pos]
			if ok {
				answered++
			}
			sc := scoring.Response{Section: sec.ID, Answered: ok, Correct: ok && r.Correct}
			if ok && sec.Scoring.Kind == spec.KindWeighted {
				o, found := l.snaps[it.Pos].Option(r.Answer)
				if !found {
					return Result{}, fmt.Errorf("%w: item %d: answer %s is not an option of the frozen question", errCorrupt, it.Pos, r.Answer)
				}
				sc.Weight = o.Weight
			}
			rs = append(rs, sc)
		}
		sr, err := scoring.ScoreSection(sec, track, rs)
		if err != nil {
			return Result{}, fmt.Errorf("%w: score section %s: %w", errCorrupt, sec.ID, err)
		}
		res.Sections = append(res.Sections, SectionResult{
			ID: sr.ID, Items: sec.Items, Answered: answered, Score: sr.Score, Max: sr.Max, Min: sr.Min, Met: sr.Met,
		})
		res.Total += sr.Score
		res.MaxTotal += sr.Max
	}
	return res, nil
}
