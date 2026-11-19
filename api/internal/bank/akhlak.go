package bank

import (
	"fmt"
	"math/rand"

	"latihseleksi/api/internal/gen"
)

// AkhlakMinReviewed is the number of reviewed AKHLAK items below which the
// section is not offered. It equals the item count of the AKHLAK section in
// the rbb-2024 spec (pinned by a test), so Create cannot fail on Pick.
const AkhlakMinReviewed = 20

// AkhlakMinPerValue is the number of reviewed items every core value needs
// before the section is offered, so a draw cannot ignore a value. The figure
// is the authors' choice, not from a source.
const AkhlakMinPerValue = 3

// AkhlakSource serves the AKHLAK section from the reviewed AKHLAK items.
type AkhlakSource struct{ *Bank }

var _ ItemSource = AkhlakSource{}

// Sections lists the AKHLAK section only once enough items are reviewed, in
// total and for every value, so exams fall back to the available sections
// instead of failing on Pick or measuring only some values.
func (a AkhlakSource) Sections() []string {
	cov := a.Bank.Coverage(SectionAKHLAK)
	total := 0
	for _, v := range a.Bank.AkhlakRubric.Values {
		n := cov[v.ID].Reviewed
		if n < AkhlakMinPerValue {
			return nil
		}
		total += n
	}
	if total < AkhlakMinReviewed {
		return nil
	}
	return []string{SectionAKHLAK}
}

// Pick returns n distinct reviewed AKHLAK items.
func (a AkhlakSource) Pick(section string, n int, seed int64) ([]gen.Draft, error) {
	if section != SectionAKHLAK {
		return nil, fmt.Errorf("bank: AKHLAK source has no section %q", section)
	}
	return a.Bank.Pick(SectionAKHLAK, n, seed)
}

// akhlakDraft is the TKP draft with an AKHLAK pattern id (akhlak-<value>) and
// the same "bobot = penilaian penyusun" step.
func (e TKPEntry) akhlakDraft(rng *rand.Rand, seed int64) gen.Draft {
	d := e.draft(rng, seed)
	d.PatternID = "akhlak-" + e.Aspek
	return d
}
