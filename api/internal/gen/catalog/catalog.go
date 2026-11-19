// Package catalog lists every registered question family.
package catalog

import (
	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/analitis"
	"latihseleksi/api/internal/gen/aritmetika"
	"latihseleksi/api/internal/gen/deretangka"
	"latihseleksi/api/internal/gen/derethuruf"
	"latihseleksi/api/internal/gen/himpunan"
	"latihseleksi/api/internal/gen/kecepatan"
	"latihseleksi/api/internal/gen/kerjasama"
	"latihseleksi/api/internal/gen/perbandingan"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/gen/ratarata"
	"latihseleksi/api/internal/gen/silogisme"
	"latihseleksi/api/internal/gen/umur"
)

// All returns the families in display order.
func All() []gen.Family {
	return []gen.Family{persentase.Family{}, kerjasama.Family{}, kecepatan.Family{}, ratarata.Family{}, deretangka.Family{}, perbandingan.Family{}, silogisme.Family{}, umur.Family{}, aritmetika.Family{}, himpunan.Family{}, derethuruf.Family{}, analitis.Family{}}
}

// Subset returns the families with the given ids, in the order given. It
// panics on an unknown id: callers are startup wiring with literal ids.
func Subset(ids ...string) []gen.Family {
	byID := map[string]gen.Family{}
	for _, f := range All() {
		byID[f.ID()] = f
	}
	out := make([]gen.Family, 0, len(ids))
	for _, id := range ids {
		f, ok := byID[id]
		if !ok {
			panic("catalog: unknown family " + id)
		}
		out = append(out, f)
	}
	return out
}
