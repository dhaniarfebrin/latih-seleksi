// Package derethuruf generates letter-series TIU questions: alphabet
// sequences with fixed steps, interleaved sequences, letters mixed with
// numbers, and growing steps (pattern TIU.numerik.deret-huruf).
package derethuruf

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/optkit"
)

const (
	PatternID   = "TIU.numerik.deret-huruf"
	maxAttempts = 200
	letters     = 26
	maxNumber   = 99
)

// Misconception ids attached to wrong options.
const (
	MisPolaLain     = "pola-lain"
	MisSelisihTetap = "selisih-tetap"
	MisLompatLagi   = "lompat-lagi"
	MisSalahHitung  = "salah-hitung-dekat"
)

// term is one member of a series: a letter (1..26) or a number.
type term struct {
	letter bool
	v      int
}

func (t term) String() string {
	if t.letter {
		return string(rune('A' + t.v - 1))
	}
	return fmt.Sprint(t.v)
}

func (t term) valid() bool {
	if t.letter {
		return t.v >= 1 && t.v <= letters
	}
	return t.v >= 1 && t.v <= maxNumber
}

// Family implements gen.Family for letter-series questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string    { return PatternID }
func (Family) Label() string { return "Deret huruf" }

func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisPolaLain:     "Melanjutkan pola dari kelompok suku yang salah pada deret berselang-seling",
		MisSelisihTetap: "Memakai selisih terakhir secara tetap padahal selisihnya terus bertambah",
		MisLompatLagi:   "Melompat sekali lagi dari suku yang sudah benar",
		MisSalahHitung:  "Salah hitung kecil: hasil mendekati jawaban benar",
	}
}

// arithmetic returns n terms starting at first with a fixed step.
func arithmetic(first term, step, n int) []term {
	out := make([]term, n)
	for i := range out {
		out[i] = term{first.letter, first.v + i*step}
	}
	return out
}

// growing returns a series whose first step is step and whose steps then
// grow by one in the same direction.
func growing(first term, step, n int) []term {
	dir := 1
	if step < 0 {
		dir = -1
	}
	out := []term{first}
	for i := 1; i < n; i++ {
		out = append(out, term{first.letter, out[i-1].v + step + dir*(i-1)})
	}
	return out
}

func interleave(a, b []term, n int) []term {
	out := make([]term, 0, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			out = append(out, a[i/2])
		} else {
			out = append(out, b[i/2])
		}
	}
	return out
}

func nonZero(rng *rand.Rand, max int) int {
	v := rng.Intn(max) + 1
	if rng.Intn(2) == 0 {
		return -v
	}
	return v
}

// predict extrapolates the series by finite differences. period is 1 or 2 and
// order the polynomial degree per subsequence (1: constant step, 2: steps that
// grow by a constant). ok is false when a subsequence has too few terms or
// does not fit, or mixes letters and numbers.
func predict(ts []term, period, order int) (term, bool) {
	sub := make([]term, 0, len(ts))
	for i := len(ts) % period; i < len(ts); i += period {
		sub = append(sub, ts[i])
	}
	// the subsequence of the next index is the one starting at len%period.
	if len(sub) < order+2 {
		return term{}, false
	}
	for _, t := range sub {
		if t.letter != sub[0].letter {
			return term{}, false
		}
	}
	diffs := make([]int, len(sub))
	for i, t := range sub {
		diffs[i] = t.v
	}
	var last []int
	for o := 0; o < order; o++ {
		last = append(last, diffs[len(diffs)-1])
		next := make([]int, len(diffs)-1)
		for i := range next {
			next[i] = diffs[i+1] - diffs[i]
		}
		diffs = next
	}
	for _, d := range diffs {
		if d != diffs[0] {
			return term{}, false
		}
	}
	v := diffs[0]
	for i := len(last) - 1; i >= 0; i-- {
		v += last[i]
	}
	return term{sub[0].letter, v}, true
}

// answer returns the unique next term under every hypothesis that fits and
// the period of the first one that fits; ok is false when none fits or two
// fitting hypotheses disagree.
func answer(ts []term) (key term, period int, ok bool) {
	for _, h := range [][2]int{{1, 1}, {1, 2}, {2, 1}} {
		p, fits := predict(ts, h[0], h[1])
		if !fits {
			continue
		}
		if period != 0 && key != p {
			return term{}, 0, false
		}
		if period == 0 {
			key, period = p, h[0]
		}
	}
	return key, period, period != 0
}

// Generate returns the question for seed and band: band 1 a constant step,
// band 2 two interleaved sequences or letters mixed with numbers, band 3
// steps that keep growing or longer interleaved sequences.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return optkit.Retry("derethuruf", seed, band, maxAttempts, func(rng *rand.Rand) (gen.Draft, bool) {
		letter := func() term { return term{true, rng.Intn(letters) + 1} }
		var ts []term
		var how string
		switch band {
		case gen.Band1:
			s := nonZero(rng, 4)
			ts = arithmetic(letter(), s, 5)
			how = fmt.Sprintf("Setiap huruf bergeser %+d urutan alfabet dari huruf sebelumnya.", s)
		case gen.Band2:
			if rng.Intn(2) == 0 {
				s1, s2 := nonZero(rng, 3), nonZero(rng, 3)
				ts = interleave(arithmetic(letter(), s1, 4), arithmetic(letter(), s2, 4), 7)
				how = fmt.Sprintf("Suku ke-1, 3, 5, ... bergeser %+d dan suku ke-2, 4, 6, ... bergeser %+d.", s1, s2)
			} else {
				s1, s2 := nonZero(rng, 3), nonZero(rng, 4)
				ts = interleave(arithmetic(letter(), s1, 4), arithmetic(term{false, rng.Intn(20) + 2}, s2, 4), 7)
				how = fmt.Sprintf("Huruf (suku ganjil) bergeser %+d dan angka (suku genap) berubah %+d.", s1, s2)
			}
		default:
			if rng.Intn(2) == 0 {
				s := rng.Intn(2) + 1
				if rng.Intn(2) == 0 {
					s = -s
				}
				ts = growing(letter(), s, 6)
				how = fmt.Sprintf("Selisih antarhuruf mulai dari %+d dan bertambah 1 selangkah demi selangkah ke arah yang sama.", s)
			} else {
				s1, s2 := nonZero(rng, 3), nonZero(rng, 3)
				ts = interleave(arithmetic(letter(), s1, 4), arithmetic(letter(), s2, 4), 8)
				how = fmt.Sprintf("Suku ganjil bergeser %+d dan suku genap bergeser %+d.", s1, s2)
			}
		}
		for _, t := range ts {
			if !t.valid() {
				return gen.Draft{}, false
			}
		}
		key, period, ok := answer(ts)
		if !ok || !key.valid() {
			return gen.Draft{}, false
		}
		names := make([]string, len(ts))
		for i, t := range ts {
			names[i] = t.String()
		}
		stem := fmt.Sprintf("Lanjutkan deret berikut: %s, ... Suku berikutnya adalah ...", strings.Join(names, ", "))
		steps := []string{how, fmt.Sprintf("Suku berikutnya adalah %s.", key)}

		var cands []optkit.Cand
		push := func(t term, mis string) {
			if t.valid() {
				cands = append(cands, optkit.Cand{Text: t.String(), Mis: mis})
			}
		}
		same := func(v int) term { return term{key.letter, v} }
		n := len(ts)
		own := ts[n-period]
		push(same(key.v+(key.v-own.v)), MisLompatLagi)
		if period == 1 {
			push(same(ts[n-1].v+(ts[n-1].v-ts[n-2].v)), MisSelisihTetap)
		} else {
			other := ts[n-1]
			push(term{other.letter, other.v + (other.v - ts[n-3].v)}, MisPolaLain)
			push(other, MisPolaLain)
		}
		for d := 1; d <= 6; d++ {
			push(same(key.v+d), MisSalahHitung)
			push(same(key.v-d), MisSalahHitung)
		}
		draft, err := optkit.Build(PatternID, stem, seed, key.String(), cands, steps)
		return draft, err == nil
	})
}
