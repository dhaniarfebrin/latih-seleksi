// Package gen defines the contract shared by all template question families.
package gen

// Band is the difficulty band of a generated question (1 = easiest).
type Band int

const (
	Band1 Band = 1
	Band2 Band = 2
	Band3 Band = 3
)

// Option is one answer choice. MisconceptionID is empty for the correct option.
// Weight is the score of a weighted (TKP) option; 0 for right-wrong sections.
type Option struct {
	Key             string
	Text            string
	MisconceptionID string
	Weight          int
}

// Draft is a generated question with its key and solution steps.
type Draft struct {
	PatternID string
	Stem      string
	Options   []Option
	AnswerKey string
	Steps     []string
	Seed      int64
}

// Family generates questions of one pattern.
type Family interface {
	ID() string
	// Label is the user-facing name of the pattern.
	Label() string
	// Misconceptions maps each misconception id the family can emit to its
	// user-facing label. Ids are unique per pattern, not globally.
	Misconceptions() map[string]string
	Generate(seed int64, band Band) (Draft, error)
}
