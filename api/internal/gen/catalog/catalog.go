// Package catalog lists every registered question family.
package catalog

import (
	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/kecepatan"
	"latihseleksi/api/internal/gen/kerjasama"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/gen/ratarata"
)

// All returns the families in display order.
func All() []gen.Family {
	return []gen.Family{persentase.Family{}, kerjasama.Family{}, kecepatan.Family{}, ratarata.Family{}}
}
