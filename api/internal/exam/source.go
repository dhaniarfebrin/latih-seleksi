package exam

import (
	"fmt"
	"math/rand"

	"latihseleksi/api/internal/bank"
	"latihseleksi/api/internal/gen"
)

// ItemSource serves ready questions per section; bank.Bank (TWK, TKP) and
// the generator adapter below both implement it.
type ItemSource = bank.ItemSource

// sectionTIU is the section served by generated families.
const sectionTIU = "TIU"

// seedTries bounds the retries for one position when a generated stem repeats.
const seedTries = 5

// GeneratedSource serves items of one section from template families.
type GeneratedSource struct {
	section  string
	families []gen.Family
}

var _ ItemSource = GeneratedSource{}

// NewGeneratedSource serves TIU from the given families (use catalog.All()).
func NewGeneratedSource(families []gen.Family) GeneratedSource {
	return NewGeneratedSourceFor(sectionTIU, families)
}

// NewGeneratedSourceFor serves the given section from the given families.
func NewGeneratedSourceFor(section string, families []gen.Family) GeneratedSource {
	return GeneratedSource{section: section, families: families}
}

func (g GeneratedSource) Sections() []string { return []string{g.section} }

// Pick returns n items, round-robin over the families starting at a seeded
// offset, each with a seeded band 1-3. It avoids repeated stems when it can;
// after seedTries repeats it keeps the last draft rather than fail. The result
// is deterministic per seed.
func (g GeneratedSource) Pick(section string, n int, seed int64) ([]gen.Draft, error) {
	if section != g.section {
		return nil, fmt.Errorf("exam: generated source has no section %q", section)
	}
	if len(g.families) == 0 {
		return nil, fmt.Errorf("exam: no families for section %s", section)
	}
	if n < 0 {
		return nil, fmt.Errorf("exam: negative item count %d", n)
	}
	rng := rand.New(rand.NewSource(seed))
	offset := rng.Intn(len(g.families))
	seen := make(map[string]bool, n)
	out := make([]gen.Draft, 0, n)
	for i := 0; i < n; i++ {
		fam := g.families[(offset+i)%len(g.families)]
		var draft gen.Draft
		var lastErr error
		got := false
		for try := 0; try < seedTries; try++ {
			d, err := fam.Generate(rng.Int63(), gen.Band(rng.Intn(3)+1))
			if err != nil {
				lastErr = err
				continue
			}
			draft, got = d, true
			if !seen[d.Stem] {
				break
			}
		}
		if !got {
			return nil, fmt.Errorf("exam: generate item %d (%s): %w", i+1, fam.ID(), lastErr)
		}
		seen[draft.Stem] = true
		out = append(out, draft)
	}
	return out, nil
}
