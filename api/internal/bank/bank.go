// Package bank holds the TWK and TKP item banks as embedded data. TWK keys are
// cited quotes checked against a verbatim corpus; TKP weights are the
// authors' assessment (draft until a human reviews them), tied to a
// behaviour guideline item. Nothing here uses a database.
package bank

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	"latihseleksi/api/internal/spec"
)

//go:embed data/twk/*.json data/tkp/*.json data/rubric/*.json data/corpus/*.txt
var dataFS embed.FS

// Section ids, matching the exam spec.
const (
	SectionTWK = "TWK"
	SectionTKP = "TKP"
)

// specID is the exam spec whose aspect lists the bank must use.
const specID = "skd-2024"

// Item status. Only StatusReviewed items are served by Pick.
const (
	StatusDraft    = "draft"
	StatusReviewed = "reviewed"
)

// TWKEntry is one TWK fact question. The key (Correct) must be contained in
// Quote, which must be contained in the corpus of SourceID.
type TWKEntry struct {
	ID          string   `json:"id"`
	Aspek       string   `json:"aspek"`
	SourceID    string   `json:"source_id"`
	Ref         string   `json:"ref"`
	Quote       string   `json:"quote"`
	Question    string   `json:"question"`
	Correct     string   `json:"correct"`
	Distractors []string `json:"distractors"`
	Status      string   `json:"status"`
	ReviewedBy  string   `json:"reviewed_by,omitempty"`
}

// Who reviewed a reviewed TWK entry. Agent review is a check against the
// source text, not a human reading.
const (
	ReviewerAgent = "agent"
	ReviewerHuman = "human"
)

// TKPOption is one TKP answer with its weight (1-5).
type TKPOption struct {
	Text   string `json:"text"`
	Weight int    `json:"weight"`
}

// RubricRef points at one behaviour item of the rubric.
type RubricRef struct {
	Value string `json:"value"`
	Butir string `json:"butir"`
}

// TKPEntry is one TKP situational question.
type TKPEntry struct {
	ID        string      `json:"id"`
	Aspek     string      `json:"aspek"`
	Scenario  string      `json:"scenario"`
	Options   []TKPOption `json:"options"`
	RubricRef RubricRef   `json:"rubric_ref"`
	Status    string      `json:"status"`
	// ReviewedBy is "human" for a reviewed entry and empty for a draft: TKP
	// weights are never agent-reviewed.
	ReviewedBy string `json:"reviewed_by,omitempty"`
}

// Butir is one behaviour item (a, b, c) of a core value.
type Butir struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// RubricValue is one core value with its behaviour items.
type RubricValue struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	Butir []Butir `json:"butir"`
}

// Rubric is the behaviour guideline copied from SE MenPANRB 20/2021.
type Rubric struct {
	Source string        `json:"source"`
	Values []RubricValue `json:"values"`
}

// Bank is the loaded and validated bank.
type Bank struct {
	TWK    []TWKEntry
	TKP    []TKPEntry
	Rubric Rubric
	Corpus map[string]string
}

// Load parses and validates every embedded file against the aspects of the
// skd-2024 spec. Any invalid entry or duplicate id is an error.
func Load() (*Bank, error) {
	reg, err := spec.Load()
	if err != nil {
		return nil, err
	}
	s, ok := reg.Get(specID)
	if !ok {
		return nil, fmt.Errorf("bank: spec %q not found", specID)
	}
	return loadFS(dataFS, aspectsOf(s, SectionTWK), aspectsOf(s, SectionTKP))
}

func aspectsOf(s spec.Spec, section string) []string {
	for _, sec := range s.Sections {
		if sec.ID == section {
			return sec.Aspects
		}
	}
	return nil
}

func loadFS(fsys fs.FS, twkAspects, tkpAspects []string) (*Bank, error) {
	b := &Bank{Corpus: map[string]string{}}

	corpus, err := fs.Glob(fsys, "data/corpus/*.txt")
	if err != nil {
		return nil, err
	}
	for _, name := range corpus {
		raw, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("bank: read %s: %w", name, err)
		}
		id := strings.TrimSuffix(name[strings.LastIndex(name, "/")+1:], ".txt")
		b.Corpus[id] = string(raw)
	}

	rubrics, err := fs.Glob(fsys, "data/rubric/*.json")
	if err != nil {
		return nil, err
	}
	if len(rubrics) != 1 {
		return nil, fmt.Errorf("bank: want exactly 1 rubric file, got %d", len(rubrics))
	}
	if err := readJSON(fsys, rubrics[0], &b.Rubric); err != nil {
		return nil, err
	}
	if err := b.Rubric.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", rubrics[0], err)
	}

	seen := map[string]string{}
	twkFiles, err := fs.Glob(fsys, "data/twk/*.json")
	if err != nil {
		return nil, err
	}
	for _, name := range twkFiles {
		var es []TWKEntry
		if err := readJSON(fsys, name, &es); err != nil {
			return nil, err
		}
		for _, e := range es {
			if err := e.Validate(b.Corpus, twkAspects); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			if prev, dup := seen[e.ID]; dup {
				return nil, fmt.Errorf("%s: duplicate id %q (also in %s)", name, e.ID, prev)
			}
			seen[e.ID] = name
			b.TWK = append(b.TWK, e)
		}
	}
	tkpFiles, err := fs.Glob(fsys, "data/tkp/*.json")
	if err != nil {
		return nil, err
	}
	for _, name := range tkpFiles {
		var es []TKPEntry
		if err := readJSON(fsys, name, &es); err != nil {
			return nil, err
		}
		for _, e := range es {
			if err := e.Validate(b.Rubric, tkpAspects); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			if prev, dup := seen[e.ID]; dup {
				return nil, fmt.Errorf("%s: duplicate id %q (also in %s)", name, e.ID, prev)
			}
			seen[e.ID] = name
			b.TKP = append(b.TKP, e)
		}
	}
	sort.Slice(b.TWK, func(i, j int) bool { return b.TWK[i].ID < b.TWK[j].ID })
	sort.Slice(b.TKP, func(i, j int) bool { return b.TKP[i].ID < b.TKP[j].ID })
	return b, nil
}

// readJSON decodes one file strictly (unknown fields are errors).
func readJSON(fsys fs.FS, name string, v any) error {
	raw, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("bank: read %s: %w", name, err)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return fmt.Errorf("%s: trailing data after JSON value", name)
	}
	return nil
}

// Coverage counts items per aspect and status.
type Coverage struct {
	Reviewed int
	Draft    int
}

// Coverage reports, for one section, the item counts per aspect.
func (b *Bank) Coverage(section string) map[string]Coverage {
	out := map[string]Coverage{}
	add := func(aspek, status string) {
		c := out[aspek]
		if status == StatusReviewed {
			c.Reviewed++
		} else {
			c.Draft++
		}
		out[aspek] = c
	}
	switch section {
	case SectionTWK:
		for _, e := range b.TWK {
			add(e.Aspek, e.Status)
		}
	case SectionTKP:
		for _, e := range b.TKP {
			add(e.Aspek, e.Status)
		}
	}
	return out
}
