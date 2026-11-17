// Package analitis generates seating-order puzzles for the TIU verbal
// section (pattern TIU.verbal.analitis). The key comes from enumerating every
// arrangement: constraints are added until exactly one arrangement is left.
package analitis

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/optkit"
)

const (
	PatternID      = "TIU.verbal.analitis"
	maxAttempts    = 40
	maxSlots       = 7 // 7! = 5040 arrangements bounds the solver cost
	maxConstraints = 9
	sampleTries    = 80
)

// Misconception ids attached to wrong options.
const (
	MisUrutanDibalik = "urutan-dibalik"
	MisGeser         = "geser-satu"
	MisKandidatLain  = "kandidat-lain"
)

var pool = []string{"Andi", "Budi", "Citra", "Dewi", "Eko", "Fajar", "Gita", "Hana", "Indra"}

type kind int

const (
	kindAdjacent    kind = iota // x next to y
	kindNotAdjacent             // x not next to y
	kindLeft                    // x somewhere left of y
	kindImmLeft                 // x directly left of y
	kindFirst                   // x at the left end
	kindLast                    // x at the right end
)

type constraint struct {
	kind kind
	x, y int // person indexes
}

// holds reports whether pos (person -> slot) satisfies c in a row of n slots.
func (c constraint) holds(pos []int, n int) bool {
	px, py := pos[c.x], pos[c.y]
	switch c.kind {
	case kindAdjacent:
		return px-py == 1 || py-px == 1
	case kindNotAdjacent:
		return px-py > 1 || py-px > 1
	case kindLeft:
		return px < py
	case kindImmLeft:
		return px+1 == py
	case kindFirst:
		return px == 0
	}
	return px == n-1
}

func (c constraint) text(names []string) string {
	x, y := names[c.x], names[c.y]
	switch c.kind {
	case kindAdjacent:
		return fmt.Sprintf("%s duduk bersebelahan dengan %s", x, y)
	case kindNotAdjacent:
		return fmt.Sprintf("%s tidak duduk bersebelahan dengan %s", x, y)
	case kindLeft:
		return fmt.Sprintf("%s duduk di sebelah kiri %s, boleh berjarak", x, y)
	case kindImmLeft:
		return fmt.Sprintf("%s duduk tepat di sebelah kiri %s", x, y)
	case kindFirst:
		return fmt.Sprintf("%s duduk paling kiri", x)
	}
	return fmt.Sprintf("%s duduk paling kanan", x)
}

var (
	permsMu sync.Mutex
	perms   = map[int][][]int{}
)

// allPerms returns every arrangement of n people (person -> slot), cached.
func allPerms(n int) [][]int {
	permsMu.Lock()
	defer permsMu.Unlock()
	if p, ok := perms[n]; ok {
		return p
	}
	var out [][]int
	var rec func(cur []int, used int)
	rec = func(cur []int, used int) {
		if len(cur) == n {
			out = append(out, append([]int(nil), cur...))
			return
		}
		for v := 0; v < n; v++ {
			if used&(1<<v) == 0 {
				rec(append(cur, v), used|1<<v)
			}
		}
	}
	rec(nil, 0)
	perms[n] = out
	return out
}

// Family implements gen.Family for seating puzzles.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string    { return PatternID }
func (Family) Label() string { return "Penalaran analitis" }

func (Family) Misconceptions() map[string]string {
	return map[string]string{
		MisUrutanDibalik: "Menghitung urutan dari sisi kanan padahal yang ditanyakan dari kiri",
		MisGeser:         "Salah menghitung satu urutan: memilih tetangga dari posisi yang benar",
		MisKandidatLain:  "Memilih jawaban yang melanggar salah satu syarat",
	}
}

func slotsFor(band gen.Band) int {
	return map[gen.Band]int{gen.Band1: 5, gen.Band2: 6, gen.Band3: 7}[band]
}

var numberWord = map[int]string{5: "Lima", 6: "Enam", 7: "Tujuh"}

// Generate returns the question for seed and band (5, 6 or 7 seats). The
// constraints are sampled from a hidden arrangement and kept only when they
// remove at least one arrangement; generation stops at exactly one survivor.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return optkit.Retry("analitis", seed, band, maxAttempts, func(rng *rand.Rand) (gen.Draft, bool) {
		n := slotsFor(band)
		names := append([]string(nil), pool...)
		rng.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })
		names = names[:n]
		truth := rng.Perm(n)
		survivors := allPerms(n)
		var cs []constraint
		for try := 0; try < sampleTries && len(survivors) > 1 && len(cs) < maxConstraints; try++ {
			c := constraint{kind: kind(rng.Intn(6)), x: rng.Intn(n), y: rng.Intn(n)}
			if (c.kind < kindFirst && c.x == c.y) || !c.holds(truth, n) {
				continue
			}
			var next [][]int
			for _, p := range survivors {
				if c.holds(p, n) {
					next = append(next, p)
				}
			}
			if len(next) == len(survivors) {
				continue
			}
			survivors, cs = next, append(cs, c)
		}
		if len(survivors) != 1 {
			return gen.Draft{}, false
		}

		// People named in a position clue are not asked about: that would
		// reveal the answer directly.
		fixed := map[int]bool{}
		for _, c := range cs {
			if c.kind == kindFirst || c.kind == kindLast {
				fixed[c.x] = true
			}
		}
		var askable []int
		for i := 0; i < n; i++ {
			if !fixed[i] {
				askable = append(askable, i)
			}
		}
		if len(askable) == 0 {
			return gen.Draft{}, false
		}
		who := askable[rng.Intn(len(askable))]
		at := make([]int, n) // slot -> person
		for person, slot := range truth {
			at[slot] = person
		}
		askName := rng.Intn(2) == 0
		slot := truth[who]
		var question string
		text := func(s int) string { return fmt.Sprint(s + 1) }
		if askName {
			question = fmt.Sprintf("Siapa yang duduk di urutan ke-%d dari kiri?", slot+1)
			text = func(s int) string { return names[at[s]] }
		} else {
			question = fmt.Sprintf("Pada urutan ke berapa dari kiri %s duduk?", names[who])
		}

		var parts []string
		for i, c := range cs {
			parts = append(parts, fmt.Sprintf("(%d) %s.", i+1, c.text(names)))
		}
		list := strings.Join(names[:n-1], ", ") + " dan " + names[n-1]
		stem := fmt.Sprintf("%s orang, yaitu %s, duduk berjajar dalam satu baris dan semuanya menghadap ke depan. Diketahui: %s %s",
			numberWord[n], list, strings.Join(parts, " "), question)
		var order []string
		for s := 0; s < n; s++ {
			order = append(order, names[at[s]])
		}
		steps := []string{
			fmt.Sprintf("Dari %d syarat, hanya satu susunan yang memenuhi semuanya.", len(cs)),
			fmt.Sprintf("Susunan dari kiri ke kanan: %s.", strings.Join(order, ", ")),
			fmt.Sprintf("Jawaban: %s.", text(slot)),
		}

		var cands []optkit.Cand
		cands = append(cands, optkit.Cand{Text: text(n - 1 - slot), Mis: MisUrutanDibalik})
		for _, s := range []int{slot - 1, slot + 1} {
			if s >= 0 && s < n {
				cands = append(cands, optkit.Cand{Text: text(s), Mis: MisGeser})
			}
		}
		for _, s := range rng.Perm(n) {
			cands = append(cands, optkit.Cand{Text: text(s), Mis: MisKandidatLain})
		}
		draft, err := optkit.Build(PatternID, stem, seed, text(slot), cands, steps)
		return draft, err == nil
	})
}
