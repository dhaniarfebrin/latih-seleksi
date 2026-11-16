package bank

import (
	"fmt"
	"math/rand"
	"sort"

	"latihseleksi/api/internal/gen"
)

// ItemSource serves ready questions per section, so drill and simulation do
// not care whether an item is generated or comes from a bank.
type ItemSource interface {
	Sections() []string
	// Pick returns n distinct items of the section, deterministic per seed.
	Pick(section string, n int, seed int64) ([]gen.Draft, error)
}

var _ ItemSource = (*Bank)(nil)

// Sections lists the sections this bank serves.
func (b *Bank) Sections() []string { return []string{SectionTWK, SectionTKP} }

var optionKeys = []string{"A", "B", "C", "D", "E"}

// Pick returns n distinct reviewed items of the section. Draft items are never
// served. It errors when fewer than n reviewed items exist.
func (b *Bank) Pick(section string, n int, seed int64) ([]gen.Draft, error) {
	var pool []func(rng *rand.Rand, seed int64) gen.Draft
	switch section {
	case SectionTWK:
		for _, e := range b.TWK {
			if e.Status == StatusReviewed {
				pool = append(pool, e.draft)
			}
		}
	case SectionTKP:
		for _, e := range b.TKP {
			if e.Status == StatusReviewed {
				pool = append(pool, e.draft)
			}
		}
	default:
		return nil, fmt.Errorf("bank: unknown section %q", section)
	}
	if n < 0 || n > len(pool) {
		return nil, fmt.Errorf("bank: want %d reviewed %s items, have %d", n, section, len(pool))
	}
	rng := rand.New(rand.NewSource(seed))
	idx := rng.Perm(len(pool))[:n]
	sort.Ints(idx)
	out := make([]gen.Draft, 0, n)
	for i, k := range idx {
		out = append(out, pool[k](rng, seed+int64(i)))
	}
	rng.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out, nil
}

func (e TWKEntry) draft(rng *rand.Rand, seed int64) gen.Draft {
	texts := append([]string{e.Correct}, e.Distractors...)
	rng.Shuffle(len(texts), func(i, j int) { texts[i], texts[j] = texts[j], texts[i] })
	d := gen.Draft{PatternID: "twk-" + e.Aspek, Stem: e.Question, Seed: seed}
	for i, t := range texts {
		o := gen.Option{Key: optionKeys[i], Text: t}
		if t == e.Correct {
			d.AnswerKey = o.Key
		} else {
			o.MisconceptionID = "distraktor"
		}
		d.Options = append(d.Options, o)
	}
	d.Steps = []string{fmt.Sprintf("%s: %q", e.Ref, e.Quote)}
	return d
}

func (e TKPEntry) draft(rng *rand.Rand, seed int64) gen.Draft {
	// Shuffle a copy so the key position is not the data order; keys and
	// steps are built after the shuffle so they stay consistent.
	opts := append([]TKPOption(nil), e.Options...)
	rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	d := gen.Draft{PatternID: "tkp-" + e.Aspek, Stem: e.Scenario, Seed: seed}
	for i, o := range opts {
		key := optionKeys[i]
		if o.Weight == 5 && d.AnswerKey == "" {
			d.AnswerKey = key
		}
		d.Options = append(d.Options, gen.Option{Key: key, Text: o.Text})
		d.Steps = append(d.Steps, fmt.Sprintf("%s: bobot %d", key, o.Weight))
	}
	d.Steps = append(d.Steps, fmt.Sprintf("Rujukan butir: %s/%s (bobot = penilaian penyusun, bukan kunci resmi)", e.RubricRef.Value, e.RubricRef.Butir))
	return d
}
