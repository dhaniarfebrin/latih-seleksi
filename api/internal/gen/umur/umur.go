// Package umur generates age-puzzle TIU questions (pattern TIU.numerik.umur).
package umur

import (
	"fmt"
	"math/rand"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/optkit"
)

const (
	PatternID   = "TIU.numerik.umur"
	maxAttempts = 50
	maxAge      = 80
)

// Misconception ids attached to wrong options.
const (
	MisTertukar       = "orang-tertukar"
	MisSetengahJumlah = "setengah-jumlah"
	MisLupaTahunLalu  = "lupa-selisih-tahun"
	MisGeserTerbalik  = "geser-terbalik"
	MisUmurLampau     = "umur-lampau"
	MisUmurMendatang  = "umur-mendatang"
	MisSalahHitung    = "salah-hitung-dekat"
)

// pair names two relatives; gap is the smallest believable age difference.
type pair struct {
	older, younger string
	gap            int
}

var pairs = []pair{{"Ayah", "Anak", 20}, {"Ibu", "Putri", 20}, {"Kakak", "Adik", 1}, {"Paman", "Keponakan", 12}, {"Kakek", "Cucu", 40}}

// Family implements gen.Family for age questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string    { return PatternID }
func (Family) Label() string { return "Soal umur" }

func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisTertukar:       "Menjawab umur orang lain, bukan yang ditanyakan",
		MisSetengahJumlah: "Membagi dua jumlah umur, padahal umur keduanya berbeda",
		MisLupaTahunLalu:  "Lupa memperhitungkan pergeseran tahun pada jumlah umur",
		MisGeserTerbalik:  "Membalik arah pergeseran tahun: menambah padahal mengurangi",
		MisUmurLampau:     "Menjawab umur pada waktu lampau, bukan umur sekarang",
		MisUmurMendatang:  "Menjawab umur pada waktu mendatang, bukan umur sekarang",
		MisSalahHitung:    "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band: band 1 sum and difference,
// band 2 a ratio with a past sum, band 3 a past ratio with a future sum.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return optkit.Retry("umur", seed, band, maxAttempts, func(rng *rand.Rand) (gen.Draft, bool) {
		p := pairs[rng.Intn(len(pairs))]
		askOlder := rng.Intn(2) == 0
		who := p.younger
		if askOlder {
			who = p.older
		}
		var (
			stem, key string
			steps     []string
			cands     []optkit.Cand
			older     int
			younger   int
		)
		ask := func(a, b int) int {
			if askOlder {
				return a
			}
			return b
		}
		txt := func(v int) string { return fmt.Sprintf("%d tahun", v) }
		switch band {
		case gen.Band1:
			b := rng.Intn(30) + 5
			d := rng.Intn(30) + 2
			a := b + d
			if a > maxAge {
				return gen.Draft{}, false
			}
			s := a + b
			stem = fmt.Sprintf("Jumlah umur %s dan %s adalah %d tahun, sedangkan selisih umur mereka %d tahun. Umur %s sekarang adalah ...", p.older, p.younger, s, d, who)
			steps = []string{
				fmt.Sprintf("Umur %s = (%d + %d) : 2 = %d tahun.", p.older, s, d, a),
				fmt.Sprintf("Umur %s = (%d - %d) : 2 = %d tahun.", p.younger, s, d, b),
				fmt.Sprintf("Yang ditanyakan adalah umur %s, yaitu %d tahun.", who, ask(a, b)),
			}
			older, younger = a, b
			k := ask(a, b)
			key = txt(k)
			cands = append(cands, optkit.Cand{Text: txt(ask(b, a)), Mis: MisTertukar})
			if s%2 == 0 {
				cands = append(cands, optkit.Cand{Text: txt(s / 2), Mis: MisSetengahJumlah})
			}
		case gen.Band2:
			b := rng.Intn(20) + 8
			k := rng.Intn(3) + 2
			n := rng.Intn(b-2) + 2
			a := k * b
			if a > maxAge {
				return gen.Draft{}, false
			}
			s := a + b - 2*n
			stem = fmt.Sprintf("Umur %s sekarang %d kali umur %s. %d tahun yang lalu, jumlah umur mereka adalah %d tahun. Umur %s sekarang adalah ...", p.older, k, p.younger, n, s, who)
			steps = []string{
				fmt.Sprintf("Misalkan umur %s sekarang x, maka umur %s sekarang %dx.", p.younger, p.older, k),
				fmt.Sprintf("(%dx - %d) + (x - %d) = %d, sehingga %dx = %d dan x = %d.", k, n, n, s, k+1, s+2*n, b),
				fmt.Sprintf("Umur %s sekarang %d tahun dan umur %s sekarang %d tahun.", p.younger, b, p.older, a),
			}
			older, younger = a, b
			key = txt(ask(a, b))
			cands = append(cands, optkit.Cand{Text: txt(ask(b, a)), Mis: MisTertukar})
			if s%(k+1) == 0 {
				x := s / (k + 1)
				cands = append(cands, optkit.Cand{Text: txt(ask(k*x, x)), Mis: MisLupaTahunLalu})
			}
			if s-2*n > 0 && (s-2*n)%(k+1) == 0 {
				x := (s - 2*n) / (k + 1)
				cands = append(cands, optkit.Cand{Text: txt(ask(k*x, x)), Mis: MisGeserTerbalik})
			}
		default:
			b := rng.Intn(23) + 8
			n := rng.Intn(b-1) + 1
			m := rng.Intn(4) + 2
			y := rng.Intn(9) + 2
			a := m*(b-n) + n
			if a > maxAge || a <= b {
				return gen.Draft{}, false
			}
			s := a + b + 2*y
			stem = fmt.Sprintf("%d tahun yang lalu, umur %s adalah %d kali umur %s. %d tahun lagi, jumlah umur mereka menjadi %d tahun. Umur %s sekarang adalah ...", n, p.older, m, p.younger, y, s, who)
			steps = []string{
				fmt.Sprintf("Umur %s = %d x (umur %s - %d) + %d.", p.older, m, p.younger, n, n),
				fmt.Sprintf("Jumlah umur sekarang = %d - 2 x %d = %d.", s, y, a+b),
				fmt.Sprintf("Dari kedua persamaan: umur %s sekarang %d tahun dan umur %s sekarang %d tahun.", p.younger, b, p.older, a),
			}
			older, younger = a, b
			key = txt(ask(a, b))
			cands = append(cands,
				optkit.Cand{Text: txt(ask(a, b) - n), Mis: MisUmurLampau},
				optkit.Cand{Text: txt(ask(a, b) + y), Mis: MisUmurMendatang},
				optkit.Cand{Text: txt(ask(b, a)), Mis: MisTertukar})
		}
		if older-younger < p.gap {
			return gen.Draft{}, false
		}
		k := ask(older, younger)
		for d := 1; d <= 6; d++ {
			if k-d > 0 {
				cands = append(cands, optkit.Cand{Text: txt(k - d), Mis: MisSalahHitung})
			}
			cands = append(cands, optkit.Cand{Text: txt(k + d), Mis: MisSalahHitung})
		}
		draft, err := optkit.Build(PatternID, stem, seed, key, cands, steps)
		return draft, err == nil
	})
}
