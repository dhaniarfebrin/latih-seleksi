// Package question holds what drill and exam share: the frozen question
// snapshot, the A-E answer check, seed generation and the common error
// sentinels. It depends only on gen.
package question

import (
	"crypto/rand"
	"encoding/binary"
	"errors"

	"latihseleksi/api/internal/gen"
)

// Sentinels shared by every feature, so the HTTP layer maps one error per meaning.
var (
	// ErrNotFound means the requested record does not exist (or is not the caller's).
	ErrNotFound = errors.New("not found")
	// ErrInvalid marks a caller mistake; wrap it with %w and a reason.
	ErrInvalid = errors.New("invalid input")
)

// answerKeys are the option letters an answer may use.
const answerKeys = "ABCDE"

// Snapshot is the frozen question stored in attempt_item.item_snapshot.
type Snapshot struct {
	Stem      string           `json:"stem"`
	Options   []SnapshotOption `json:"options"`
	AnswerKey string           `json:"answer_key"`
	Steps     []string         `json:"steps"`
	// FactID is the bank entry id, frozen so spaced repetition can key on it.
	FactID string `json:"fact_id,omitempty"`
}

// SnapshotOption is one frozen answer option.
type SnapshotOption struct {
	Key             string `json:"key"`
	Text            string `json:"text"`
	MisconceptionID string `json:"misconception_id,omitempty"`
	Weight          int    `json:"weight,omitempty"`
}

// FromDraft freezes a generated or banked draft.
func FromDraft(d gen.Draft) Snapshot {
	sn := Snapshot{Stem: d.Stem, AnswerKey: d.AnswerKey, Steps: d.Steps, FactID: d.FactID}
	for _, o := range d.Options {
		sn.Options = append(sn.Options, SnapshotOption{Key: o.Key, Text: o.Text, MisconceptionID: o.MisconceptionID, Weight: o.Weight})
	}
	return sn
}

// Option returns the option with the given key.
func (s Snapshot) Option(key string) (SnapshotOption, bool) {
	for _, o := range s.Options {
		if o.Key == key {
			return o, true
		}
	}
	return SnapshotOption{}, false
}

// ValidAnswer reports whether a is exactly one letter A-E.
func ValidAnswer(a string) bool {
	return len(a) == 1 && a[0] >= 'A' && a[0] <= answerKeys[len(answerKeys)-1]
}

// RandomSeed returns a non-negative random seed for a generator.
func RandomSeed() (int64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b[:]) >> 1), nil
}
