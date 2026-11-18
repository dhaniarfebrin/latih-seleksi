// Package himpunan generates set-counting TIU questions with two or three
// sets (pattern TIU.numerik.himpunan).
package himpunan

import (
	"fmt"
	"math/rand"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/optkit"
)

const (
	PatternID   = "TIU.numerik.himpunan"
	maxAttempts = 50
)

// Misconception ids attached to wrong options.
const (
	MisJumlahLangsung = "jumlah-langsung"
	MisLupaIrisan     = "lupa-irisan"
	MisTertukar       = "pertanyaan-tertukar"
	MisSalahHitung    = "salah-hitung-dekat"
)

var (
	nouns      = []string{"siswa", "pegawai", "peserta", "warga"}
	activities = []string{"sepak bola", "basket", "voli", "renang", "catur", "bulu tangkis"}
)

// Family implements gen.Family for set-counting questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string    { return PatternID }
func (Family) Label() string { return "Himpunan" }

func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisJumlahLangsung: "Menjumlahkan semua kelompok tanpa mengurangi anggota yang terhitung dua kali",
		MisLupaIrisan:     "Lupa menyesuaikan irisan: ada bagian yang terhitung dua kali atau terkurangi berlebihan",
		MisTertukar:       "Menjawab besaran lain, bukan yang ditanyakan (misalnya gabungan, bagian di luar, atau satu himpunan saja)",
		MisSalahHitung:    "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// Generate returns the question for seed and band. Region counts are chosen
// first, so every stem is consistent by construction: band 1 two sets with the
// overlap given, band 2 two sets with the outside count given, band 3 three
// sets with every pairwise and triple overlap given.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return optkit.Retry("himpunan", seed, band, maxAttempts, func(rng *rand.Rand) (gen.Draft, bool) {
		noun := nouns[rng.Intn(len(nouns))]
		perm := rng.Perm(len(activities))
		x, y, z := activities[perm[0]], activities[perm[1]], activities[perm[2]]
		r := func(lo, hi int) int { return lo + rng.Intn(hi-lo+1) }
		var stem string
		var key int
		var steps []string
		var cands []int
		var mis []string
		add := func(v int, m string) { cands, mis = append(cands, v), append(mis, m) }
		if band != gen.Band3 {
			onlyA, both, onlyB, none := r(2, 15), r(1, 10), r(2, 15), r(1, 12)
			a, b := onlyA+both, onlyB+both
			n := onlyA + both + onlyB + none
			if band == gen.Band1 {
				stem = fmt.Sprintf("Dari %d %s, %d gemar %s, %d gemar %s, dan %d gemar keduanya. ", n, noun, a, x, b, y, both)
				union := a + b - both
				switch rng.Intn(3) {
				case 0:
					stem += fmt.Sprintf("Banyak %s yang tidak gemar satu pun adalah ...", noun)
					key = n - union
					steps = []string{fmt.Sprintf("Gemar paling sedikit satu = %d + %d - %d = %d.", a, b, both, union), fmt.Sprintf("Tidak gemar satu pun = %d - %d = %d.", n, union, key)}
					add(n-a-b, MisLupaIrisan)
					add(union, MisTertukar)
					add(n-union+both, MisLupaIrisan)
				case 1:
					stem += fmt.Sprintf("Banyak %s yang gemar paling sedikit satu adalah ...", noun)
					key = union
					steps = []string{fmt.Sprintf("Gemar paling sedikit satu = %d + %d - %d = %d.", a, b, both, union)}
					add(a+b, MisJumlahLangsung)
					add(n-union, MisTertukar)
					add(a+b-2*both, MisLupaIrisan)
				default:
					stem += fmt.Sprintf("Banyak %s yang hanya gemar %s adalah ...", noun, x)
					key = onlyA
					steps = []string{fmt.Sprintf("Hanya gemar %s = %d - %d = %d.", x, a, both, key)}
					add(a, MisLupaIrisan)
					add(onlyB, MisTertukar)
					add(a+both, MisLupaIrisan)
				}
			} else {
				stem = fmt.Sprintf("Dari %d %s, %d gemar %s, %d gemar %s, dan %d tidak gemar keduanya. Banyak %s yang gemar keduanya adalah ...", n, noun, a, x, b, y, none, noun)
				key = both
				steps = []string{fmt.Sprintf("Gemar paling sedikit satu = %d - %d = %d.", n, none, n-none), fmt.Sprintf("Gemar keduanya = %d + %d - %d = %d.", a, b, n-none, key)}
				add(a+b-n, MisLupaIrisan)
				add(a+b, MisJumlahLangsung)
				add(n-none-a, MisTertukar)
				add(n-none-b, MisTertukar)
			}
		} else {
			oa, ob, oc := r(2, 9), r(2, 9), r(2, 9)
			ab, ac, bc, abc, none := r(1, 6), r(1, 6), r(1, 6), r(1, 6), r(1, 10)
			a, b, c := oa+ab+ac+abc, ob+ab+bc+abc, oc+ac+bc+abc
			pab, pac, pbc := ab+abc, ac+abc, bc+abc
			union := oa + ob + oc + ab + ac + bc + abc
			n := union + none
			stem = fmt.Sprintf("Dari %d %s, %d gemar %s, %d gemar %s, dan %d gemar %s. Sebanyak %d gemar %s dan %s, %d gemar %s dan %s, %d gemar %s dan %s, serta %d gemar ketiganya. ",
				n, noun, a, x, b, y, c, z, pab, x, y, pac, x, z, pbc, y, z, abc)
			sum := a + b + c
			pairs := pab + pac + pbc
			switch rng.Intn(3) {
			case 0:
				stem += fmt.Sprintf("Banyak %s yang gemar paling sedikit satu adalah ...", noun)
				key = union
				steps = []string{fmt.Sprintf("Gabungan = %d - (%d) + %d = %d.", sum, pairs, abc, union)}
				add(sum, MisJumlahLangsung)
				add(sum-pairs, MisLupaIrisan)
				add(n-union, MisTertukar)
			case 1:
				stem += fmt.Sprintf("Banyak %s yang tidak gemar satu pun adalah ...", noun)
				key = none
				steps = []string{fmt.Sprintf("Gabungan = %d - (%d) + %d = %d.", sum, pairs, abc, union), fmt.Sprintf("Tidak gemar satu pun = %d - %d = %d.", n, union, none)}
				add(union, MisTertukar)
				add(n-(sum-pairs), MisLupaIrisan)
				add(n-sum, MisJumlahLangsung)
			default:
				stem += fmt.Sprintf("Banyak %s yang gemar tepat satu dari ketiganya adalah ...", noun)
				key = oa + ob + oc
				steps = []string{fmt.Sprintf("Tepat satu = %d - 2 x (%d) + 3 x %d = %d.", sum, pairs, abc, key)}
				add(sum-pairs, MisLupaIrisan)
				add(sum-2*pairs, MisLupaIrisan)
				add(union, MisTertukar)
			}
		}
		var cs []optkit.Cand
		for i, v := range cands {
			if v >= 0 && v != key {
				cs = append(cs, optkit.Cand{Text: fmt.Sprint(v), Mis: mis[i]})
			}
		}
		for d := 1; d <= 6; d++ {
			for _, v := range []int{key + d, key - d} {
				if v >= 0 {
					cs = append(cs, optkit.Cand{Text: fmt.Sprint(v), Mis: MisSalahHitung})
				}
			}
		}
		draft, err := optkit.Build(PatternID, stem, seed, fmt.Sprint(key), cs, steps)
		return draft, err == nil
	})
}
