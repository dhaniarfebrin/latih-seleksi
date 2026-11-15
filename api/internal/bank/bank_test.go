package bank

import (
	"regexp"
	"strings"
	"testing"
)

var ws = regexp.MustCompile(`\s+`)

// squash is an independent normaliser (regexp, not the package's norm).
func squash(s string) string { return strings.ToLower(ws.ReplaceAllString(strings.TrimSpace(s), " ")) }

func loadBank(t *testing.T) *Bank {
	t.Helper()
	b, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return b
}

func TestEmbeddedDataLoads(t *testing.T) {
	b := loadBank(t)
	if len(b.TWK) < 10 || len(b.TKP) < 10 {
		t.Fatalf("seed too small: %d TWK, %d TKP", len(b.TWK), len(b.TKP))
	}
	ids := map[string]bool{}
	for _, e := range b.TWK {
		ids[e.ID] = true
	}
	for _, e := range b.TKP {
		if ids[e.ID] {
			t.Errorf("id %q used in both TWK and TKP", e.ID)
		}
		ids[e.ID] = true
	}
	if len(ids) != len(b.TWK)+len(b.TKP) {
		t.Errorf("ids not unique")
	}
}

func TestTWKKeysMatchCorpusIndependently(t *testing.T) {
	b := loadBank(t)
	for _, e := range b.TWK {
		corpus := squash(b.Corpus[e.SourceID])
		if corpus == "" {
			t.Errorf("%s: empty corpus %q", e.ID, e.SourceID)
			continue
		}
		q := squash(e.Quote)
		if !strings.Contains(corpus, q) {
			t.Errorf("%s: quote not in corpus", e.ID)
		}
		if !strings.Contains(q, squash(e.Correct)) {
			t.Errorf("%s: correct not in quote", e.ID)
		}
		for _, d := range e.Distractors {
			if strings.Contains(q, squash(d)) {
				t.Errorf("%s: distractor %q in quote", e.ID, d)
			}
		}
	}
}

func TestSeedTKPStaysDraft(t *testing.T) {
	for _, e := range loadBank(t).TKP {
		if e.Status != StatusDraft {
			t.Errorf("%s: seed TKP weights need a human, status must be draft, got %q", e.ID, e.Status)
		}
	}
}

func TestCoverage(t *testing.T) {
	b := loadBank(t)
	c := b.Coverage(SectionTWK)
	total := 0
	for _, v := range c {
		total += v.Reviewed + v.Draft
	}
	if total != len(b.TWK) || len(c) != 5 {
		t.Errorf("TWK coverage %v does not match %d entries over 5 aspects", c, len(b.TWK))
	}
}

var (
	testCorpus = map[string]string{"src": "Pasal 1\nBahasa Negara ialah\nBahasa Indonesia."}
	testAspect = []string{"bahasa-negara"}
)

func validTWK() TWKEntry {
	return TWKEntry{ID: "x", Aspek: "bahasa-negara", SourceID: "src", Ref: "Pasal 1",
		Quote: "Bahasa Negara ialah Bahasa Indonesia.", Question: "Bahasa negara ialah ...",
		Correct: "Bahasa Indonesia", Distractors: []string{"Bahasa Jawa", "Bahasa Melayu", "Bahasa Inggris"}, Status: StatusReviewed}
}

func TestTWKValidateRejects(t *testing.T) {
	if err := validTWK().Validate(testCorpus, testAspect); err != nil {
		t.Fatalf("valid entry rejected (line-wrapped quote must match): %v", err)
	}
	cases := map[string]func(*TWKEntry){
		"quote not in corpus":  func(e *TWKEntry) { e.Quote = "Bahasa Negara ialah Bahasa Jawa." },
		"correct not in quote": func(e *TWKEntry) { e.Correct = "Bahasa Jawa"; e.Distractors[0] = "Bahasa Bali" },
		"distractor in quote":  func(e *TWKEntry) { e.Distractors[0] = "Bahasa Negara" },
		"no ref":               func(e *TWKEntry) { e.Ref = "" },
		"bad aspek":            func(e *TWKEntry) { e.Aspek = "lain" },
		"unknown source":       func(e *TWKEntry) { e.SourceID = "none" },
		"duplicate option":     func(e *TWKEntry) { e.Distractors[1] = e.Distractors[0] },
		"too few distractors":  func(e *TWKEntry) { e.Distractors = e.Distractors[:2] },
		"bad status":           func(e *TWKEntry) { e.Status = "ok" },
	}
	for name, mut := range cases {
		e := validTWK()
		e.Distractors = append([]string(nil), e.Distractors...)
		mut(&e)
		if err := e.Validate(testCorpus, testAspect); err == nil {
			t.Errorf("%s: want error, got nil", name)
		}
	}
}

func testRubric() Rubric {
	r := Rubric{}
	for _, id := range []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7"} {
		r.Values = append(r.Values, RubricValue{ID: id, Label: id, Butir: []Butir{{"a", "x"}, {"b", "y"}, {"c", "z"}}})
	}
	return r
}

func validTKP() TKPEntry {
	return TKPEntry{ID: "k", Aspek: "tik", Scenario: "s", RubricRef: RubricRef{Value: "a1", Butir: "a"}, Status: StatusReviewed,
		Options: []TKPOption{{"o1", 5}, {"o2", 4}, {"o3", 3}, {"o4", 2}, {"o5", 1}}}
}

func TestTKPValidateRejects(t *testing.T) {
	if err := validTKP().Validate(testRubric(), []string{"tik"}); err != nil {
		t.Fatalf("valid entry rejected: %v", err)
	}
	cases := map[string]func(*TKPEntry){
		"no weight 5":   func(e *TKPEntry) { e.Options[0].Weight = 4 },
		"no low weight": func(e *TKPEntry) { e.Options[3].Weight = 3; e.Options[4].Weight = 3 },
		"weight 6":      func(e *TKPEntry) { e.Options[1].Weight = 6 },
		"weight 0":      func(e *TKPEntry) { e.Options[1].Weight = 0 },
		"four options":  func(e *TKPEntry) { e.Options = e.Options[:4] },
		"same text":     func(e *TKPEntry) { e.Options[1].Text = "O1" },
		"unknown value": func(e *TKPEntry) { e.RubricRef.Value = "zz" },
		"unknown butir": func(e *TKPEntry) { e.RubricRef.Butir = "z" },
		"bad aspek":     func(e *TKPEntry) { e.Aspek = "lain" },
	}
	for name, mut := range cases {
		e := validTKP()
		e.Options = append([]TKPOption(nil), e.Options...)
		mut(&e)
		if err := e.Validate(testRubric(), []string{"tik"}); err == nil {
			t.Errorf("%s: want error, got nil", name)
		}
	}
}

func TestPick(t *testing.T) {
	b := loadBank(t)
	// Fixtures: reviewed TKP (seed ones are draft) and one TWK left as draft.
	for i := range b.TKP[:6] {
		b.TKP[i].Status = StatusReviewed
	}
	b.TWK[0].Status = StatusDraft
	draftTWK := b.TWK[0].Question

	for _, section := range []string{SectionTWK, SectionTKP} {
		a, err := b.Pick(section, 5, 7)
		if err != nil {
			t.Fatalf("%s: %v", section, err)
		}
		c, _ := b.Pick(section, 5, 7)
		if len(a) != 5 {
			t.Fatalf("%s: got %d items", section, len(a))
		}
		stems := map[string]bool{}
		for i := range a {
			if a[i].Stem != c[i].Stem || a[i].AnswerKey != c[i].AnswerKey {
				t.Errorf("%s: not deterministic for the same seed", section)
			}
			if stems[a[i].Stem] {
				t.Errorf("%s: repeated item %q", section, a[i].Stem)
			}
			stems[a[i].Stem] = true
			if a[i].Stem == draftTWK {
				t.Errorf("%s: draft item served", section)
			}
			if a[i].AnswerKey == "" || len(a[i].Options) < 4 {
				t.Errorf("%s: malformed draft %+v", section, a[i])
			}
		}
	}
	if _, err := b.Pick(SectionTKP, 7, 1); err == nil {
		t.Error("asking for more reviewed items than exist must error")
	}
	if _, err := b.Pick("TIU", 1, 1); err == nil {
		t.Error("unknown section must error")
	}
}

func TestPickAnswerKeys(t *testing.T) {
	b := loadBank(t)
	for i := range b.TKP {
		b.TKP[i].Status = StatusReviewed
	}
	byStem := map[string]TKPEntry{}
	for _, e := range b.TKP {
		byStem[e.Scenario] = e
	}
	twk := map[string]TWKEntry{}
	for _, e := range b.TWK {
		twk[e.Question] = e
	}
	tkpDrafts, _ := b.Pick(SectionTKP, len(b.TKP), 3)
	for _, d := range tkpDrafts {
		e := byStem[d.Stem]
		for i, o := range d.Options {
			if o.Key == d.AnswerKey && e.Options[i].Weight != 5 {
				t.Errorf("%s: key %s has weight %d", e.ID, o.Key, e.Options[i].Weight)
			}
		}
	}
	twkDrafts, _ := b.Pick(SectionTWK, len(b.TWK), 3)
	for _, d := range twkDrafts {
		e := twk[d.Stem]
		for _, o := range d.Options {
			if (o.Key == d.AnswerKey) != (o.Text == e.Correct) {
				t.Errorf("%s: key %s does not point at the correct text", e.ID, d.AnswerKey)
			}
		}
	}
}
