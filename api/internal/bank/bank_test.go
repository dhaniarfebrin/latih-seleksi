package bank

import (
	"fmt"
	"math/rand"
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

func TestTKPReviewedOnlyByHuman(t *testing.T) {
	for _, e := range loadBank(t).TKP {
		if e.Status == StatusReviewed && e.ReviewedBy != ReviewerHuman {
			t.Errorf("%s: TKP weights need a human, reviewed_by must be %q, got %q", e.ID, ReviewerHuman, e.ReviewedBy)
		}
	}
}

func TestTKPBankShape(t *testing.T) {
	b := loadBank(t)
	if len(b.TKP) != 45 {
		t.Errorf("want 45 TKP items, got %d", len(b.TKP))
	}
	perAspek := map[string]int{}
	pos5 := map[int]int{}
	for _, e := range b.TKP {
		perAspek[e.Aspek]++
		w := map[int]bool{}
		for i, o := range e.Options {
			w[o.Weight] = true
			if o.Weight == 5 {
				pos5[i]++
			}
		}
		if len(w) != 5 {
			t.Errorf("%s: weights must be 1..5 once each, got %v", e.ID, e.Options)
		}
	}
	for a, n := range perAspek {
		if n < 7 || n > 8 {
			t.Errorf("aspek %s: want 7-8 items, got %d", a, n)
		}
	}
	if len(pos5) < 3 {
		t.Errorf("weight-5 position too predictable in data: %v", pos5)
	}
}

func TestTKPDraftShufflesAndKeepsKey(t *testing.T) {
	e := validTKP()
	keys := map[string]bool{}
	for seed := int64(0); seed < 40; seed++ {
		d := e.draft(rand.New(rand.NewSource(seed)), seed)
		if len(d.Options) != 5 || len(d.Steps) != 6 {
			t.Fatalf("want 5 options and 6 steps, got %d and %d", len(d.Options), len(d.Steps))
		}
		weight := map[string]int{}
		for _, o := range e.Options {
			weight[o.Text] = o.Weight
		}
		for i, o := range d.Options {
			if want := fmt.Sprintf("%s: bobot %d", o.Key, weight[o.Text]); d.Steps[i] != want {
				t.Fatalf("seed %d: step %q does not match option %+v", seed, d.Steps[i], o)
			}
			if o.Weight != weight[o.Text] {
				t.Fatalf("seed %d: option %+v weight != %d", seed, o, weight[o.Text])
			}
			if o.Key == d.AnswerKey && weight[o.Text] != 5 {
				t.Fatalf("seed %d: answer key %s is not the weight-5 option", seed, d.AnswerKey)
			}
		}
		keys[d.AnswerKey] = true
	}
	if len(keys) < 3 {
		t.Errorf("answer key position too predictable across seeds: %v", keys)
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

func TestTWKHasEnoughReviewedPerAspek(t *testing.T) {
	const min = 15
	for a, v := range loadBank(t).Coverage(SectionTWK) {
		if v.Reviewed < min {
			t.Errorf("aspek %s: want >= %d reviewed TWK items, got %d", a, min, v.Reviewed)
		}
	}
}

var (
	testCorpus = map[string]string{"src": "Pasal 1\nBahasa Negara ialah\nBahasa Indonesia."}
	testAspect = []string{"bahasa-negara"}
)

func validTWK() TWKEntry {
	return TWKEntry{ID: "x", Aspek: "bahasa-negara", SourceID: "src", Ref: "Pasal 1",
		Quote: "Bahasa Negara ialah Bahasa Indonesia.", Question: "Bahasa negara ialah ...",
		Correct: "Bahasa Indonesia", Distractors: []string{"Bahasa Jawa", "Bahasa Melayu", "Bahasa Inggris"}, Status: StatusReviewed, ReviewedBy: ReviewerAgent}
}

// blockCorpus has two ayat in one pasal and a second pasal.
var blockCorpus = map[string]string{"src": "Pasal 2\n(1) Warna bendera ialah merah dan putih.\n(2) Ukuran bendera ialah dua pertiga dari panjang.\n\nPasal 3\nBendera dikibarkan setiap hari."}

func blockTWK() TWKEntry {
	return TWKEntry{ID: "b", Aspek: "bahasa-negara", SourceID: "src", Ref: "Pasal 2 ayat (1)",
		Quote: "Warna bendera ialah merah dan putih.", Question: "Warna bendera ialah ...",
		Correct: "merah dan putih", Distractors: []string{"hijau dan kuning", "biru dan hitam", "ungu dan jingga"}, Status: StatusReviewed, ReviewedBy: ReviewerAgent}
}

func TestTWKValidateRejectsDistractorFromSamePasal(t *testing.T) {
	if err := blockTWK().Validate(blockCorpus, testAspect); err != nil {
		t.Fatalf("valid entry rejected: %v", err)
	}
	e := blockTWK()
	e.Distractors = []string{"hijau dan kuning", "dua pertiga dari panjang", "ungu dan jingga"}
	if err := e.Validate(blockCorpus, testAspect); err == nil {
		t.Error("distractor taken from ayat (2) of the same pasal must be rejected")
	}
	e = blockTWK()
	e.Distractors = []string{"hijau dan kuning", "dikibarkan setiap hari", "ungu dan jingga"}
	if err := e.Validate(blockCorpus, testAspect); err != nil {
		t.Errorf("text from another pasal is not a same-pasal conflict: %v", err)
	}
}

func TestTWKValidateRejects(t *testing.T) {
	if err := validTWK().Validate(testCorpus, testAspect); err != nil {
		t.Fatalf("valid entry rejected (line-wrapped quote must match): %v", err)
	}
	cases := map[string]func(*TWKEntry){
		"quote not in corpus":   func(e *TWKEntry) { e.Quote = "Bahasa Negara ialah Bahasa Jawa." },
		"correct not in quote":  func(e *TWKEntry) { e.Correct = "Bahasa Jawa"; e.Distractors[0] = "Bahasa Bali" },
		"distractor in quote":   func(e *TWKEntry) { e.Distractors[0] = "Bahasa Negara" },
		"no ref":                func(e *TWKEntry) { e.Ref = "" },
		"bad aspek":             func(e *TWKEntry) { e.Aspek = "lain" },
		"unknown source":        func(e *TWKEntry) { e.SourceID = "none" },
		"duplicate option":      func(e *TWKEntry) { e.Distractors[1] = e.Distractors[0] },
		"too few distractors":   func(e *TWKEntry) { e.Distractors = e.Distractors[:2] },
		"bad status":            func(e *TWKEntry) { e.Status = "ok" },
		"reviewed no reviewer":  func(e *TWKEntry) { e.ReviewedBy = "" },
		"reviewed bad reviewer": func(e *TWKEntry) { e.ReviewedBy = "bot" },
		"draft with reviewer":   func(e *TWKEntry) { e.Status = StatusDraft },
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
	return TKPEntry{ID: "k", Aspek: "tik", Scenario: "s", RubricRef: RubricRef{Value: "a1", Butir: "a"}, Status: StatusReviewed, ReviewedBy: ReviewerHuman,
		Options: []TKPOption{{"o1", 5}, {"o2", 4}, {"o3", 3}, {"o4", 2}, {"o5", 1}}}
}

func TestTKPValidateRejects(t *testing.T) {
	if err := validTKP().Validate(testRubric(), []string{"tik"}); err != nil {
		t.Fatalf("valid entry rejected: %v", err)
	}
	cases := map[string]func(*TKPEntry){
		"no weight 5":               func(e *TKPEntry) { e.Options[0].Weight = 4 },
		"no low weight":             func(e *TKPEntry) { e.Options[3].Weight = 3; e.Options[4].Weight = 3 },
		"weight 6":                  func(e *TKPEntry) { e.Options[1].Weight = 6 },
		"weight 0":                  func(e *TKPEntry) { e.Options[1].Weight = 0 },
		"four options":              func(e *TKPEntry) { e.Options = e.Options[:4] },
		"same text":                 func(e *TKPEntry) { e.Options[1].Text = "O1" },
		"unknown value":             func(e *TKPEntry) { e.RubricRef.Value = "zz" },
		"unknown butir":             func(e *TKPEntry) { e.RubricRef.Butir = "z" },
		"bad aspek":                 func(e *TKPEntry) { e.Aspek = "lain" },
		"reviewed without reviewer": func(e *TKPEntry) { e.ReviewedBy = "" },
		"reviewed by agent":         func(e *TKPEntry) { e.ReviewedBy = ReviewerAgent },
		"draft with reviewer":       func(e *TKPEntry) { e.Status = StatusDraft },
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
	// Fixture: one TWK item left as draft.
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
	if _, err := b.Pick(SectionTKP, len(b.TKP)+1, 1); err == nil {
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
		weight := map[string]int{}
		for _, o := range e.Options {
			weight[o.Text] = o.Weight
		}
		for _, o := range d.Options {
			if o.Weight != weight[o.Text] {
				t.Fatalf("%s: option %+v weight != %d", e.ID, o, weight[o.Text])
			}
			if o.Key == d.AnswerKey && weight[o.Text] != 5 {
				t.Errorf("%s: key %s has weight %d", e.ID, o.Key, weight[o.Text])
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

func TestDrawTWK(t *testing.T) {
	b := loadBank(t)
	ids := b.PatternIDs()
	if len(ids) != 5 {
		t.Fatalf("pattern ids = %v, want 5 aspects", ids)
	}
	byID := map[string]TWKEntry{}
	pool := map[string]int{}
	for _, e := range b.TWK {
		byID[e.ID] = e
		if e.Status == StatusReviewed {
			pool[DrillPrefix+e.Aspek]++
		}
	}
	for _, id := range ids {
		if !strings.HasPrefix(id, "TWK.") {
			t.Errorf("%s: want TWK. prefix", id)
		}
		if l, ok := b.PatternLabel(id); !ok || l == "" {
			t.Errorf("%s: no label", id)
		}
		got, err := b.Draw(id, 5, 42)
		if err != nil {
			t.Fatal(err)
		}
		again, _ := b.Draw(id, 5, 42)
		if len(got) != min(5, pool[id]) {
			t.Errorf("%s: drew %d, pool %d", id, len(got), pool[id])
		}
		seen := map[string]bool{}
		for i, d := range got {
			e, ok := byID[d.FactID]
			if !ok || e.Status != StatusReviewed || DrillPrefix+e.Aspek != id {
				t.Errorf("%s: bad fact %q", id, d.FactID)
			}
			if seen[d.FactID] {
				t.Errorf("%s: duplicate %s", id, d.FactID)
			}
			seen[d.FactID] = true
			if d.Stem != again[i].Stem || d.AnswerKey != again[i].AnswerKey {
				t.Errorf("%s: not deterministic", id)
			}
			var key string
			for _, o := range d.Options {
				if o.Key == d.AnswerKey {
					key = o.Text
				}
			}
			if !strings.Contains(squash(d.Steps[0]), squash(key)) {
				t.Errorf("%s: key %q not in quote step", d.FactID, key)
			}
		}
	}
	all, _ := b.Draw(ids[0], 1000, 1)
	if len(all) != pool[ids[0]] {
		t.Errorf("small pool: drew %d, want whole pool %d", len(all), pool[ids[0]])
	}
	if _, err := b.Draw("TWK.nope", 3, 1); err == nil {
		t.Error("unknown pattern accepted")
	}
}
