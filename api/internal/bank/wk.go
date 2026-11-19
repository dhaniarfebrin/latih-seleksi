package bank

import (
	"fmt"

	"latihseleksi/api/internal/gen"
)

// SectionWK is the BUMN Wawasan Kebangsaan section. It draws the same
// reviewed fact items as TWK, so its answer keys stay cited quotes.
const SectionWK = "WK"

// WKSource serves the WK section from the reviewed TWK entries of the bank.
type WKSource struct{ *Bank }

var _ ItemSource = WKSource{}

// Sections lists the only section this source serves.
func (WKSource) Sections() []string { return []string{SectionWK} }

// Pick returns n distinct reviewed TWK items. Each keeps its PatternID
// (twk-<aspek>) and FactID so reviews and the dashboard stay consistent.
func (w WKSource) Pick(section string, n int, seed int64) ([]gen.Draft, error) {
	if section != SectionWK {
		return nil, fmt.Errorf("bank: WK source has no section %q", section)
	}
	return w.Bank.Pick(SectionTWK, n, seed)
}
