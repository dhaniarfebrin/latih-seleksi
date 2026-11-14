// Package drill creates, answers and scores drill attempts.
package drill

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/store"
)

const (
	ModeDrill     = "drill"
	DefaultN      = 8
	DefaultBand   = 2
	maxN          = 20
	maxTimeMS     = 24 * 60 * 60 * 1000
	seedRetries   = 5
	answerKeyChar = "ABCDE"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrInvalid   = errors.New("invalid input")
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

// snapshot is the frozen question stored in attempt_item.item_snapshot.
type snapshot struct {
	Stem      string           `json:"stem"`
	Options   []snapshotOption `json:"options"`
	AnswerKey string           `json:"answer_key"`
	Steps     []string         `json:"steps"`
}

type snapshotOption struct {
	Key             string `json:"key"`
	Text            string `json:"text"`
	MisconceptionID string `json:"misconception_id,omitempty"`
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
	store    *store.Store
	ordered  []gen.Family // catalog order, for listing
	families map[string]gen.Family
}

func New(s *store.Store) *Service {
	ordered := catalog.All()
	families := make(map[string]gen.Family, len(ordered))
	for _, f := range ordered {
		families[f.ID()] = f
	}
	return &Service{store: s, ordered: ordered, families: families}
}

// Patterns lists the registered families in catalog order.
func (s *Service) Patterns() []Pattern {
	out := make([]Pattern, 0, len(s.ordered))
	for _, f := range s.ordered {
		out = append(out, Pattern{ID: f.ID(), Label: f.Label(), Misconceptions: f.Misconceptions()})
	}
	return out
}

func randomSeed() (int64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b[:]) >> 1), nil
}

// Create generates n questions with random seeds and freezes them.
func (s *Service) Create(ctx context.Context, patternID string, n, band int) (AttemptView, error) {
	fam, ok := s.families[patternID]
	if !ok {
		return AttemptView{}, invalid("unknown pattern_id %q", patternID)
	}
	if n < 1 || n > maxN {
		return AttemptView{}, invalid("n must be between 1 and %d", maxN)
	}
	if band < 1 || band > 3 {
		return AttemptView{}, invalid("band must be between 1 and 3")
	}

	items := make([]store.Item, 0, n)
	for pos := 1; pos <= n; pos++ {
		var draft gen.Draft
		var err error
		for try := 0; try < seedRetries; try++ {
			var seed int64
			if seed, err = randomSeed(); err != nil {
				return AttemptView{}, err
			}
			if draft, err = fam.Generate(seed, gen.Band(band)); err == nil {
				break
			}
		}
		if err != nil {
			return AttemptView{}, fmt.Errorf("generate item %d: %w", pos, err)
		}
		raw, err := json.Marshal(toSnapshot(draft))
		if err != nil {
			return AttemptView{}, err
		}
		items = append(items, store.Item{Pos: pos, Seed: draft.Seed, Snapshot: raw})
	}

	a, err := s.store.CreateAttempt(ctx, store.Attempt{Mode: ModeDrill, PatternID: patternID, Band: band, N: n}, items)
	if err != nil {
		return AttemptView{}, err
	}
	return s.Get(ctx, a.ID)
}

func toSnapshot(d gen.Draft) snapshot {
	sn := snapshot{Stem: d.Stem, AnswerKey: d.AnswerKey, Steps: d.Steps}
	for _, o := range d.Options {
		sn.Options = append(sn.Options, snapshotOption{Key: o.Key, Text: o.Text, MisconceptionID: o.MisconceptionID})
	}
	return sn
}

func (s *Service) load(ctx context.Context, id string) (store.Attempt, map[int]snapshot, []store.Item, []store.Response, error) {
	a, items, resps, err := s.store.GetAttempt(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return a, nil, nil, nil, ErrNotFound
	}
	if err != nil {
		return a, nil, nil, nil, err
	}
	snaps := make(map[int]snapshot, len(items))
	for _, it := range items {
		var sn snapshot
		if err := json.Unmarshal(it.Snapshot, &sn); err != nil {
			return a, nil, nil, nil, fmt.Errorf("decode snapshot %d: %w", it.Pos, err)
		}
		snaps[it.Pos] = sn
	}
	return a, snaps, items, resps, nil
}

func responseView(r store.Response, sn snapshot) *ResponseView {
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
func (s *Service) Get(ctx context.Context, id string) (AttemptView, error) {
	a, snaps, items, resps, err := s.load(ctx, id)
	if err != nil {
		return AttemptView{}, err
	}
	byPos := make(map[int]store.Response, len(resps))
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
func (s *Service) Answer(ctx context.Context, id string, pos int, answer string, timeMS int) (ResponseView, error) {
	if len(answer) != 1 || !containsByte(answerKeyChar, answer[0]) {
		return ResponseView{}, invalid("answer must be one of A-E")
	}
	if timeMS < 0 || timeMS > maxTimeMS {
		return ResponseView{}, invalid("time_ms must be between 0 and %d", maxTimeMS)
	}
	a, snaps, _, resps, err := s.load(ctx, id)
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
	chosen, found := snapshotOption{}, false
	for _, o := range sn.Options {
		if o.Key == answer {
			chosen, found = o, true
		}
	}
	if !found {
		return ResponseView{}, invalid("answer %s is not an option of this question", answer)
	}
	stored, err := s.store.SaveResponse(ctx, id, store.Response{
		Pos: pos, Answer: answer, Correct: answer == sn.AnswerKey, TimeMS: timeMS, MisconceptionID: chosen.MisconceptionID,
	})
	if err != nil {
		return ResponseView{}, err
	}
	return *responseView(stored, sn), nil
}

func containsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

// Submit closes the attempt and returns score, mean time and misconception tally.
// Unanswered positions count as incorrect. Repeating the call is safe.
func (s *Service) Submit(ctx context.Context, id string) (Result, error) {
	submittedAt, err := s.store.MarkSubmitted(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return Result{}, ErrNotFound
	}
	if err != nil {
		return Result{}, err
	}
	a, _, resps, err := s.store.GetAttempt(ctx, id)
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
	return res, nil
}
