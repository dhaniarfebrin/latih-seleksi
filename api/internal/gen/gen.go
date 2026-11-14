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
type Option struct {
	Key             string
	Text            string
	MisconceptionID string
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

// Params are the numbers behind one question. For chained percentages,
// Changes holds the signed percent of each step (positive = increase) and
// Context selects the sentence template.
type Params struct {
	Changes []int
	Context int
}

// Answer is the exact result of Solve: the net change in percent is Num/Den.
type Answer struct {
	Num  int64
	Den  int64
	Text string
}

// Family generates questions of one pattern and solves them independently.
type Family interface {
	ID() string
	Generate(seed int64, band Band) (Draft, error)
	Solve(p Params) (Answer, []string)
}
