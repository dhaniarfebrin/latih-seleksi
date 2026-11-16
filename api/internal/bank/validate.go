package bank

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// heading matches a corpus line that opens a block ("Pasal 25", "Angka 4").
var heading = regexp.MustCompile(`^(Pasal|Angka) \S+$`)

// block returns the normalised text of the Pasal/Angka block of the corpus
// that contains quote, or "" when no block does. A corpus without headings is
// one block.
func block(text, quote string) string {
	var blocks []string
	var cur []string
	for _, line := range strings.Split(text, "\n") {
		if heading.MatchString(strings.TrimSpace(line)) && len(cur) > 0 {
			blocks = append(blocks, strings.Join(cur, " "))
			cur = nil
		}
		cur = append(cur, line)
	}
	blocks = append(blocks, strings.Join(cur, " "))
	for _, b := range blocks {
		if nb := norm(b); strings.Contains(nb, quote) {
			return nb
		}
	}
	return ""
}

// norm lowercases and collapses whitespace so quotes match the corpus
// regardless of line wrapping.
func norm(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

func validStatus(s string) bool { return s == StatusDraft || s == StatusReviewed }

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// Validate checks one TWK entry against the corpus and the allowed aspects.
// All problems are reported together.
func (e TWKEntry) Validate(corpus map[string]string, aspects []string) error {
	var errs []error
	fail := func(format string, a ...any) {
		errs = append(errs, fmt.Errorf("twk %q: "+format, append([]any{e.ID}, a...)...))
	}
	if strings.TrimSpace(e.ID) == "" {
		fail("id is required")
	}
	if !contains(aspects, e.Aspek) {
		fail("aspek %q is not one of %v", e.Aspek, aspects)
	}
	if strings.TrimSpace(e.Ref) == "" {
		fail("ref is required")
	}
	if strings.TrimSpace(e.Question) == "" {
		fail("question is required")
	}
	if !validStatus(e.Status) {
		fail("status %q must be draft or reviewed", e.Status)
	}
	switch {
	case e.Status == StatusReviewed && e.ReviewedBy != ReviewerAgent && e.ReviewedBy != ReviewerHuman:
		fail("reviewed_by %q must be %q or %q for a reviewed entry", e.ReviewedBy, ReviewerAgent, ReviewerHuman)
	case e.Status == StatusDraft && e.ReviewedBy != "":
		fail("reviewed_by %q must be empty for a draft entry", e.ReviewedBy)
	}
	text, ok := corpus[e.SourceID]
	if !ok {
		fail("source_id %q has no corpus file", e.SourceID)
	}
	quote, correct := norm(e.Quote), norm(e.Correct)
	if quote == "" {
		fail("quote is required")
	} else if ok && !strings.Contains(norm(text), quote) {
		fail("quote is not found in corpus %q", e.SourceID)
	}
	if correct == "" {
		fail("correct is required")
	} else if !strings.Contains(quote, correct) {
		fail("correct %q is not contained in the quote", e.Correct)
	}
	if n := len(e.Distractors); n < 3 || n > 4 {
		fail("want 3-4 distractors, got %d", n)
	}
	// A distractor must be wrong in the whole pasal, not only in the quote.
	whole := quote
	if ok && quote != "" {
		if b := block(text, quote); b != "" {
			whole = b
		}
	}
	seen := map[string]bool{correct: true}
	for _, d := range e.Distractors {
		nd := norm(d)
		switch {
		case nd == "":
			fail("empty distractor")
		case strings.Contains(whole, nd):
			fail("distractor %q appears in the same pasal as the quote", d)
		case seen[nd]:
			fail("duplicate option %q", d)
		}
		seen[nd] = true
	}
	return errors.Join(errs...)
}

// Validate checks one TKP entry against the rubric and the allowed aspects.
func (e TKPEntry) Validate(r Rubric, aspects []string) error {
	var errs []error
	fail := func(format string, a ...any) {
		errs = append(errs, fmt.Errorf("tkp %q: "+format, append([]any{e.ID}, a...)...))
	}
	if strings.TrimSpace(e.ID) == "" {
		fail("id is required")
	}
	if !contains(aspects, e.Aspek) {
		fail("aspek %q is not one of %v", e.Aspek, aspects)
	}
	if strings.TrimSpace(e.Scenario) == "" {
		fail("scenario is required")
	}
	if !validStatus(e.Status) {
		fail("status %q must be draft or reviewed", e.Status)
	}
	switch {
	case e.Status == StatusReviewed && e.ReviewedBy != ReviewerHuman:
		fail("reviewed_by %q must be %q for a reviewed entry (TKP weights are never agent-reviewed)", e.ReviewedBy, ReviewerHuman)
	case e.Status == StatusDraft && e.ReviewedBy != "":
		fail("reviewed_by %q must be empty for a draft entry", e.ReviewedBy)
	}
	if len(e.Options) != 5 {
		fail("want exactly 5 options, got %d", len(e.Options))
	}
	var has5, hasLow bool
	seen := map[string]bool{}
	for _, o := range e.Options {
		if o.Weight < 1 || o.Weight > 5 {
			fail("weight %d out of range 1-5", o.Weight)
		}
		has5 = has5 || o.Weight == 5
		hasLow = hasLow || (o.Weight >= 1 && o.Weight <= 2)
		n := norm(o.Text)
		if n == "" {
			fail("empty option text")
		} else if seen[n] {
			fail("duplicate option %q", o.Text)
		}
		seen[n] = true
	}
	if !has5 {
		fail("no option has weight 5")
	}
	if !hasLow {
		fail("no option has weight <= 2")
	}
	if !r.has(e.RubricRef) {
		fail("rubric_ref %+v is not in the rubric", e.RubricRef)
	}
	return errors.Join(errs...)
}

func (r Rubric) has(ref RubricRef) bool {
	for _, v := range r.Values {
		if v.ID != ref.Value {
			continue
		}
		for _, b := range v.Butir {
			if b.ID == ref.Butir {
				return true
			}
		}
	}
	return false
}

func (r Rubric) validate() error {
	if len(r.Values) != 7 {
		return fmt.Errorf("rubric: want 7 values, got %d", len(r.Values))
	}
	ids := map[string]bool{}
	for _, v := range r.Values {
		if v.ID == "" || ids[v.ID] {
			return fmt.Errorf("rubric: empty or duplicate value id %q", v.ID)
		}
		ids[v.ID] = true
		if len(v.Butir) != 3 {
			return fmt.Errorf("rubric: value %q wants 3 butir, got %d", v.ID, len(v.Butir))
		}
		for _, b := range v.Butir {
			if b.ID == "" || strings.TrimSpace(b.Text) == "" {
				return fmt.Errorf("rubric: value %q has an empty butir", v.ID)
			}
		}
	}
	return nil
}
