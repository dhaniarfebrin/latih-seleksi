package silogisme

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
)

// ---- parsing the question text ----------------------------------------------

var (
	stemRe    = regexp.MustCompile(`^Premis 1: (.+)\. Premis 2: (.+)\. Kesimpulan yang tepat adalah \.\.\.$`)
	allRe     = regexp.MustCompile(`^Semua (.+) adalah (.+)$`)
	noRe      = regexp.MustCompile(`^Tidak ada (.+) yang merupakan (.+)$`)
	someRe    = regexp.MustCompile(`^Sebagian (.+) adalah (.+)$`)
	someNotRe = regexp.MustCompile(`^Sebagian (.+) bukan (.+)$`)
)

// tstmt is a categorical statement read from text: typ is A, E, I or O.
type tstmt struct {
	typ  byte
	x, y string
}

func parseCat(text string) (tstmt, bool) {
	for _, r := range []struct {
		re  *regexp.Regexp
		typ byte
	}{{allRe, 'A'}, {noRe, 'E'}, {someRe, 'I'}, {someNotRe, 'O'}} {
		if m := r.re.FindStringSubmatch(text); m != nil {
			return tstmt{r.typ, m[1], m[2]}, true
		}
	}
	return tstmt{}, false
}

func splitStem(t *testing.T, stem string) (string, string) {
	t.Helper()
	m := stemRe.FindStringSubmatch(stem)
	if m == nil {
		t.Fatalf("unrecognised stem %q", stem)
	}
	return m[1], m[2]
}

// ---- oracle 1: textbook table of valid syllogisms ----------------------------

// Moods are major premise, minor premise and conclusion types, by figure.
// Fifteen forms are valid without existential import; nine more need it.
var (
	tableValid = map[int][]string{
		1: {"AAA", "EAE", "AII", "EIO"},
		2: {"EAE", "AEE", "EIO", "AOO"},
		3: {"IAI", "AII", "OAO", "EIO"},
		4: {"AEE", "IAI", "EIO"},
	}
	tableImport = map[int][]string{
		1: {"AAI", "EAO"},
		2: {"EAO", "AEO"},
		3: {"AAI", "EAO"},
		4: {"AAI", "EAO", "AEO"},
	}
)

func inList(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// classifyEnd finds the figure and mood of a conclusion whose subject S and
// predicate P are the two non-middle terms of the premises. ok is false when
// the conclusion does not relate those two terms.
func classifyEnd(prem []tstmt, c tstmt) (fig int, mood string, ok bool) {
	terms := func(s tstmt) []string { return []string{s.x, s.y} }
	var mid string
	for _, a := range terms(prem[0]) {
		for _, b := range terms(prem[1]) {
			if a == b {
				mid = a
			}
		}
	}
	if mid == "" || c.x == mid || c.y == mid || c.x == c.y {
		return 0, "", false
	}
	var major, minor *tstmt
	for i := range prem {
		switch {
		case prem[i].x == c.y || prem[i].y == c.y:
			major = &prem[i]
		case prem[i].x == c.x || prem[i].y == c.x:
			minor = &prem[i]
		}
	}
	if major == nil || minor == nil || major == minor {
		return 0, "", false
	}
	switch majorMSubject, minorMSubject := major.x == mid, minor.x == mid; {
	case majorMSubject && !minorMSubject:
		fig = 1
	case !majorMSubject && !minorMSubject:
		fig = 2
	case majorMSubject && minorMSubject:
		fig = 3
	default:
		fig = 4
	}
	return fig, string([]byte{major.typ, minor.typ, c.typ}), true
}

// ---- oracle 2: brute force over small domains of individuals ------------------

const bruteN = 5

type verdict struct {
	valid         bool // true in every domain of bruteN individuals that satisfies the premises
	refuted       bool // false in some domain that satisfies the premises with all terms non-empty
	refutedStrong bool // the same when "Sebagian" premises also state the opposite "sebagian"
}

var bruteCache = map[string]verdict{}

// evalInd evaluates a statement on individuals; ind[i] is the membership mask
// of individual i over the three terms (bit k set: member of term k).
func evalInd(typ byte, x, y int, ind *[bruteN]int) bool {
	for _, m := range ind {
		inX, inY := m>>x&1 == 1, m>>y&1 == 1
		switch typ {
		case 'A':
			if inX && !inY {
				return false
			}
		case 'E':
			if inX && inY {
				return false
			}
		case 'I':
			if inX && inY {
				return true
			}
		case 'O':
			if inX && !inY {
				return true
			}
		}
	}
	return typ == 'A' || typ == 'E'
}

func bruteCheck(t *testing.T, prem []tstmt, c tstmt) verdict {
	t.Helper()
	idx := map[string]int{}
	for _, s := range prem {
		for _, n := range []string{s.x, s.y} {
			if _, ok := idx[n]; !ok {
				idx[n] = len(idx)
			}
		}
	}
	cx, okx := idx[c.x]
	cy, oky := idx[c.y]
	if len(idx) != 3 || !okx || !oky {
		t.Fatalf("terms do not line up: %v / %v", prem, c)
	}
	type cs struct {
		typ  byte
		x, y int
	}
	var ps []cs
	var strong []cs
	for _, s := range prem {
		ps = append(ps, cs{s.typ, idx[s.x], idx[s.y]})
		strong = append(strong, cs{s.typ, idx[s.x], idx[s.y]})
		switch s.typ {
		case 'I':
			strong = append(strong, cs{'O', idx[s.x], idx[s.y]})
		case 'O':
			strong = append(strong, cs{'I', idx[s.x], idx[s.y]})
		}
	}
	key := fmt.Sprint(ps, cs{c.typ, cx, cy})
	if v, ok := bruteCache[key]; ok {
		return v
	}
	v := verdict{valid: true}
	sat := false
	total := 1
	for i := 0; i < bruteN; i++ {
		total *= 8
	}
	for code := 0; code < total; code++ {
		var ind [bruteN]int
		rest := code
		cover := 0
		for i := range ind {
			ind[i] = rest % 8
			rest /= 8
			cover |= ind[i]
		}
		holds := func(list []cs) bool {
			for _, s := range list {
				if !evalInd(s.typ, s.x, s.y, &ind) {
					return false
				}
			}
			return true
		}
		if !holds(ps) {
			continue
		}
		sat = true
		cOK := evalInd(c.typ, cx, cy, &ind)
		if !cOK {
			v.valid = false
		}
		if cover == 7 && !cOK {
			v.refuted = true
			if holds(strong) {
				v.refutedStrong = true
			}
		}
	}
	v.valid = v.valid && sat
	bruteCache[key] = v
	return v
}

// ---- oracle for conditional questions ----------------------------------------

type tlit struct {
	id  int
	neg bool
}

func (l tlit) not() tlit { return tlit{l.id, !l.neg} }

type tprop struct {
	cond bool
	a, b tlit
}

var litByText = func() map[string]tlit {
	m := map[string]tlit{}
	for si, sc := range scenarios {
		for pi, pr := range sc {
			m[pr.pos] = tlit{si*100 + pi, false}
			m[pr.neg] = tlit{si*100 + pi, true}
		}
	}
	return m
}()

func lowerFirst(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func parseCond(t *testing.T, text string) tprop {
	t.Helper()
	lookup := func(s string) tlit {
		l, ok := litByText[s]
		if !ok {
			t.Fatalf("unknown proposition %q in %q", s, text)
		}
		return l
	}
	if rest, ok := strings.CutPrefix(text, "Jika "); ok {
		a, b, found := strings.Cut(rest, " maka ")
		if !found {
			t.Fatalf("conditional without maka: %q", text)
		}
		return tprop{true, lookup(a), lookup(b)}
	}
	return tprop{false, lookup(lowerFirst(text)), tlit{}}
}

func ttrue(p tprop, w map[int]bool) bool {
	val := func(l tlit) bool { return w[l.id] != l.neg }
	if !p.cond {
		return val(p.a)
	}
	return !val(p.a) || val(p.b)
}

// truthTable returns (valid, refuted, satisfiable) by listing all assignments.
func truthTable(prem []tprop, c tprop) (valid, refuted, sat bool) {
	seen := map[int]bool{}
	var ids []int
	add := func(id int) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	for _, p := range append(append([]tprop(nil), prem...), c) {
		add(p.a.id)
		if p.cond {
			add(p.b.id)
		}
	}
	sort.Ints(ids)
	valid = true
	for mask := 0; mask < 1<<len(ids); mask++ {
		w := map[int]bool{}
		for i, id := range ids {
			w[id] = mask>>i&1 == 1
		}
		ok := true
		for _, p := range prem {
			ok = ok && ttrue(p, w)
		}
		if !ok {
			continue
		}
		sat = true
		if !ttrue(c, w) {
			valid, refuted = false, true
		}
	}
	return valid && sat, refuted, sat
}

// ruleValid decides by inference rules: contrapositive, modus ponens (fact
// chaining) and hypothetical syllogism (chaining without facts).
func ruleValid(prem []tprop, c tprop) bool {
	var imps [][2]tlit
	var facts []tlit
	for _, p := range prem {
		if p.cond {
			imps = append(imps, [2]tlit{p.a, p.b}, [2]tlit{p.b.not(), p.a.not()})
		} else {
			facts = append(facts, p.a)
		}
	}
	closure := func(start []tlit) map[tlit]bool {
		got := map[tlit]bool{}
		queue := append([]tlit(nil), start...)
		for len(queue) > 0 {
			l := queue[0]
			queue = queue[1:]
			if got[l] {
				continue
			}
			got[l] = true
			for _, im := range imps {
				if im[0] == l {
					queue = append(queue, im[1])
				}
			}
		}
		return got
	}
	contradictory := func(g map[tlit]bool) bool {
		for l := range g {
			if g[l.not()] {
				return true
			}
		}
		return false
	}
	if !c.cond {
		g := closure(facts)
		return g[c.a] || contradictory(g)
	}
	g := closure(append(append([]tlit(nil), facts...), c.a))
	return g[c.b] || contradictory(g)
}

// ---- properties over generated questions ------------------------------------

func TestPoolsAreParseSafe(t *testing.T) {
	seen := map[string]bool{}
	for _, g := range termGroups {
		for _, n := range g {
			if seen[n] {
				t.Errorf("term %q repeated", n)
			}
			seen[n] = true
			for _, bad := range []string{" adalah ", " bukan ", " yang merupakan ", " maka ", "Jika", "Semua", "Sebagian", "Tidak ada"} {
				if strings.Contains(n, bad) {
					t.Errorf("term %q contains %q", n, bad)
				}
			}
		}
		if len(g) < 3 {
			t.Errorf("group too small: %v", g)
		}
	}
	for si, sc := range scenarios {
		if len(sc) < 3 {
			t.Errorf("scenario %d too small", si)
		}
		for _, pr := range sc {
			for _, s := range []string{pr.pos, pr.neg} {
				if seen[s] {
					t.Errorf("proposition %q repeated", s)
				}
				seen[s] = true
				if s != strings.ToLower(s[:1])+s[1:] || strings.Contains(s, " maka ") || strings.HasPrefix(s, "Jika") {
					t.Errorf("proposition %q is not parse safe", s)
				}
			}
			if pr.pos == pr.neg {
				t.Errorf("proposition %q equals its negation", pr.pos)
			}
		}
	}
}

// checkCategorical verifies one generated question against both oracles and
// returns the figure and mood of the key.
func checkCategorical(t *testing.T, d gen.Draft, band gen.Band, seed int64) (int, string) {
	t.Helper()
	s1, s2 := splitStem(t, d.Stem)
	var prem []tstmt
	for _, s := range []string{s1, s2} {
		p, ok := parseCat(s)
		if !ok {
			t.Fatalf("band %d seed %d: unreadable premise %q", band, seed, s)
		}
		prem = append(prem, p)
	}
	if band == gen.Band1 && (strings.Contains(d.Stem, "Sebagian")) {
		t.Fatalf("band 1 seed %d has a Sebagian premise: %q", seed, d.Stem)
	}
	if band == gen.Band2 && !strings.Contains(d.Stem, "Sebagian") {
		t.Fatalf("band 2 seed %d has no Sebagian premise: %q", seed, d.Stem)
	}
	validCount := 0
	var keyFig int
	var keyMood string
	for _, o := range d.Options {
		c, ok := parseCat(o.Text)
		if !ok {
			t.Fatalf("band %d seed %d: unreadable option %q", band, seed, o.Text)
		}
		bv := bruteCheck(t, prem, c)
		fig, mood, related := classifyEnd(prem, c)
		isKey := o.Key == d.AnswerKey
		if related {
			table := inList(tableValid[fig], mood)
			withImport := table || inList(tableImport[fig], mood)
			if table != bv.valid {
				t.Fatalf("band %d seed %d: table says valid=%v, brute force valid=%v for %q in %q", band, seed, table, bv.valid, o.Text, d.Stem)
			}
			if isKey {
				keyFig, keyMood = fig, mood
				if !table {
					t.Fatalf("band %d seed %d: key %q is not a textbook valid form (fig %d mood %s)", band, seed, o.Text, fig, mood)
				}
			} else if withImport {
				t.Fatalf("band %d seed %d: distractor %q is valid with existential import (fig %d mood %s) in %q", band, seed, o.Text, fig, mood, d.Stem)
			}
		} else if isKey {
			t.Fatalf("band %d seed %d: key %q does not relate the two end terms in %q", band, seed, o.Text, d.Stem)
		}
		if bv.valid {
			validCount++
			if !isKey {
				t.Fatalf("band %d seed %d: valid option %q is not the key in %q", band, seed, o.Text, d.Stem)
			}
			continue
		}
		if isKey {
			t.Fatalf("band %d seed %d: key %q is not valid in %q", band, seed, o.Text, d.Stem)
		}
		if !bv.refuted || !bv.refutedStrong {
			t.Fatalf("band %d seed %d: distractor %q is not refuted robustly (%+v) in %q", band, seed, o.Text, bv, d.Stem)
		}
	}
	if validCount != 1 {
		t.Fatalf("band %d seed %d: %d valid options in %q", band, seed, validCount, d.Stem)
	}
	return keyFig, keyMood
}

func checkConditional(t *testing.T, d gen.Draft, seed int64) {
	t.Helper()
	s1, s2 := splitStem(t, d.Stem)
	prem := []tprop{parseCond(t, s1), parseCond(t, s2)}
	validCount := 0
	for _, o := range d.Options {
		c := parseCond(t, o.Text)
		valid, refuted, sat := truthTable(prem, c)
		if !sat {
			t.Fatalf("seed %d: premises are inconsistent in %q", seed, d.Stem)
		}
		if rv := ruleValid(prem, c); rv != valid {
			t.Fatalf("seed %d: truth table valid=%v but rules valid=%v for %q in %q", seed, valid, rv, o.Text, d.Stem)
		}
		isKey := o.Key == d.AnswerKey
		if valid {
			validCount++
		}
		if valid != isKey {
			t.Fatalf("seed %d: option %q valid=%v but key=%v in %q", seed, o.Text, valid, isKey, d.Stem)
		}
		if !isKey && !refuted {
			t.Fatalf("seed %d: distractor %q is not refuted in %q", seed, o.Text, d.Stem)
		}
	}
	if validCount != 1 {
		t.Fatalf("seed %d: %d valid options in %q", seed, validCount, d.Stem)
	}
}

func TestPropertyKeyMatchesIndependentOracles(t *testing.T) {
	fam := Family{}
	forms := map[string]bool{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			if len(d.Options) != optionCount {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			texts := map[string]bool{}
			keys := 0
			for i, o := range d.Options {
				if o.Key != string(rune('A'+i)) {
					t.Fatalf("band %d seed %d: option %d has key %q", band, seed, i, o.Key)
				}
				if texts[o.Text] {
					t.Fatalf("band %d seed %d: duplicate option %q", band, seed, o.Text)
				}
				texts[o.Text] = true
				if o.Key == d.AnswerKey {
					keys++
					if o.MisconceptionID != "" {
						t.Fatalf("band %d seed %d: key option has misconception", band, seed)
					}
				} else if o.MisconceptionID == "" {
					t.Fatalf("band %d seed %d: distractor %s has no misconception", band, seed, o.Key)
				}
			}
			if keys != 1 {
				t.Fatalf("band %d seed %d: %d key options", band, seed, keys)
			}
			if d.Seed != seed || d.PatternID != PatternID || len(d.Steps) == 0 {
				t.Fatalf("band %d seed %d: bad metadata %d %s %v", band, seed, d.Seed, d.PatternID, d.Steps)
			}
			if band == gen.Band3 {
				checkConditional(t, d, seed)
			} else {
				fig, mood := checkCategorical(t, d, band, seed)
				forms[fmt.Sprintf("%d%s", fig, mood)] = true
			}
		}
	}
	for fig, moods := range tableValid {
		for _, m := range moods {
			if !forms[fmt.Sprintf("%d%s", fig, m)] {
				t.Errorf("valid form figure %d mood %s never appeared as a key", fig, m)
			}
		}
	}
}

// figurePremises builds premises independently of shape.premises, from the
// textbook definition of the four figures: S=0, M=1, P=2.
func figurePremises(fig int, major, minor stmtType) []stmt {
	mp, pm := [2]int{1, 2}, [2]int{2, 1}
	sm, ms := [2]int{0, 1}, [2]int{1, 0}
	var a, b [2]int
	switch fig {
	case 1:
		a, b = mp, sm
	case 2:
		a, b = pm, sm
	case 3:
		a, b = mp, ms
	default:
		a, b = pm, ms
	}
	return []stmt{{major, a[0], a[1]}, {minor, b[0], b[1]}}
}

func TestCheckerMatchesTextbookTable(t *testing.T) {
	letters := "AEIO"
	for fig := 1; fig <= 4; fig++ {
		for major := typA; major <= typO; major++ {
			for minor := typA; minor <= typO; minor++ {
				prem := figurePremises(fig, major, minor)
				for c := typA; c <= typO; c++ {
					mood := string([]byte{letters[major], letters[minor], letters[c]})
					concl := stmt{c, roleS, roleP}
					wantValid := inList(tableValid[fig], mood)
					wantImport := wantValid || inList(tableImport[fig], mood)
					if got := catValid(prem, concl); got != wantValid {
						t.Errorf("figure %d mood %s: catValid = %v, table says %v", fig, mood, got, wantValid)
					}
					// with all terms non-empty a conclusion is valid when it cannot be refuted
					gotImport := catConsistent(prem) && !catRefuted(prem, concl)
					if gotImport != wantImport {
						t.Errorf("figure %d mood %s: valid with import = %v, table says %v", fig, mood, gotImport, wantImport)
					}
				}
			}
		}
	}
}

func TestShapesPerBand(t *testing.T) {
	letters := "AEIO"
	for band, shapes := range catShapes {
		if len(shapes) == 0 {
			t.Errorf("band %d has no shapes", band)
		}
		for _, sh := range shapes {
			particular := sh.major >= typI || sh.minor >= typI
			if (band == gen.Band2) != particular {
				t.Errorf("band %d shape %+v: particular premise = %v", band, sh, particular)
			}
			p := params{kind: kindCat, names: [3]string{"guru", "dokter", "petani"}, prem: sh.premises()}
			if _, ok := solve(p); !ok {
				t.Errorf("shape fig %d %c%c cannot be solved", sh.fig, letters[sh.major], letters[sh.minor])
			}
		}
	}
	// band 1: Barbara, Celarent, Cesare, Camestres, Camenes and the fourth-figure
	// Barbara whose conclusion runs from P to S
	if len(catShapes[gen.Band1]) != 6 || len(catShapes[gen.Band2]) != 16 || len(catShapes) != 2 {
		t.Errorf("shapes by band: %d / %d", len(catShapes[gen.Band1]), len(catShapes[gen.Band2]))
	}
}

func TestCondCheckerRules(t *testing.T) {
	a, b, c := lit{0, false}, lit{1, false}, lit{2, false}
	cases := []struct {
		name  string
		prem  []prop
		concl prop
		want  bool
	}{
		{"modus ponens", []prop{imp(a, b), fact(a)}, fact(b), true},
		{"modus tollens", []prop{imp(a, b), fact(b.not())}, fact(a.not()), true},
		{"silogisme hipotetis", []prop{imp(a, b), imp(b, c)}, imp(a, c), true},
		{"kontraposisi hipotetis", []prop{imp(a, b), imp(b, c)}, imp(c.not(), a.not()), true},
		{"afirmasi konsekuen", []prop{imp(a, b), fact(b)}, fact(a), false},
		{"sangkal anteseden", []prop{imp(a, b), fact(a.not())}, fact(b.not()), false},
		{"konversi", []prop{imp(a, b), imp(b, c)}, imp(c, a), false},
		{"invers", []prop{imp(a, b), imp(b, c)}, imp(a.not(), c.not()), false},
		{"fakta tunggal dari rantai", []prop{imp(a, b), imp(b, c)}, fact(c), false},
		{"premis kontradiktif", []prop{fact(a), fact(a.not())}, fact(b), false},
	}
	for _, tc := range cases {
		if got := condValid(tc.prem, tc.concl); got != tc.want {
			t.Errorf("%s: condValid = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// ---- fixtures ---------------------------------------------------------------

func candTexts(res result) map[string]string {
	m := map[string]string{}
	for _, c := range res.cands {
		m[c.text] = c.mis
	}
	return m
}

func TestCategoricalFixtures(t *testing.T) {
	// S = guru, M = dokter, P = petani
	cases := []struct {
		name         string
		fig          int
		major, minor stmtType
		stem, key    string
		wantCands    map[string]string
		absent       []string
	}{
		{
			name: "Barbara", fig: 1, major: typA, minor: typA,
			stem: "Premis 1: Semua dokter adalah petani. Premis 2: Semua guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:  "Semua guru adalah petani",
			wantCands: map[string]string{
				"Semua petani adalah guru":             MisKonversiKeliru,
				"Tidak ada guru yang merupakan petani": MisNegasiKeliru,
				"Tidak ada petani yang merupakan guru": MisNegasiKeliru,
				"Sebagian guru bukan petani":           MisNegasiKeliru,
				"Sebagian petani bukan guru":           MisNegasiKeliru,
				"Semua dokter adalah guru":             MisKesimpulanTakTerkait,
				"Semua petani adalah dokter":           MisKesimpulanTakTerkait,
				"Tidak ada guru yang merupakan dokter": MisKesimpulanTakTerkait,
				"Sebagian guru bukan dokter":           MisKesimpulanTakTerkait,
				"Sebagian dokter bukan guru":           MisKesimpulanTakTerkait,
				"Sebagian dokter bukan petani":         MisKesimpulanTakTerkait,
			},
			// a statement that is true only if "semua" implies "sebagian" must not appear
			absent: []string{"Sebagian guru adalah petani", "Sebagian petani adalah guru", "Semua guru adalah dokter", "Sebagian guru adalah dokter"},
		},
		{name: "Celarent", fig: 1, major: typE, minor: typA,
			stem:   "Premis 1: Tidak ada dokter yang merupakan petani. Premis 2: Semua guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:    "Tidak ada guru yang merupakan petani",
			absent: []string{"Tidak ada petani yang merupakan guru", "Sebagian guru bukan petani"}},
		{name: "Darii", fig: 1, major: typA, minor: typI,
			stem: "Premis 1: Semua dokter adalah petani. Premis 2: Sebagian guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru adalah petani",
			wantCands: map[string]string{
				"Semua guru adalah petani":             MisSebagianJadiSemua,
				"Tidak ada guru yang merupakan petani": MisNegasiKeliru,
			},
			absent: []string{"Sebagian petani adalah guru"}},
		{name: "Ferio", fig: 1, major: typE, minor: typI,
			stem: "Premis 1: Tidak ada dokter yang merupakan petani. Premis 2: Sebagian guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru bukan petani",
			wantCands: map[string]string{
				"Tidak ada guru yang merupakan petani": MisSebagianJadiSemua,
				"Sebagian guru adalah petani":          MisNegasiKeliru,
			}},
		{name: "Cesare", fig: 2, major: typE, minor: typA,
			stem: "Premis 1: Tidak ada petani yang merupakan dokter. Premis 2: Semua guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:  "Tidak ada guru yang merupakan petani"},
		{name: "Camestres", fig: 2, major: typA, minor: typE,
			stem: "Premis 1: Semua petani adalah dokter. Premis 2: Tidak ada guru yang merupakan dokter. Kesimpulan yang tepat adalah ...",
			key:  "Tidak ada guru yang merupakan petani"},
		{name: "Festino", fig: 2, major: typE, minor: typI,
			stem: "Premis 1: Tidak ada petani yang merupakan dokter. Premis 2: Sebagian guru adalah dokter. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru bukan petani"},
		{name: "Baroco", fig: 2, major: typA, minor: typO,
			stem: "Premis 1: Semua petani adalah dokter. Premis 2: Sebagian guru bukan dokter. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru bukan petani"},
		{name: "Disamis", fig: 3, major: typI, minor: typA,
			stem: "Premis 1: Sebagian dokter adalah petani. Premis 2: Semua dokter adalah guru. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru adalah petani"},
		{name: "Datisi", fig: 3, major: typA, minor: typI,
			stem: "Premis 1: Semua dokter adalah petani. Premis 2: Sebagian dokter adalah guru. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru adalah petani"},
		{name: "Bocardo", fig: 3, major: typO, minor: typA,
			stem: "Premis 1: Sebagian dokter bukan petani. Premis 2: Semua dokter adalah guru. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru bukan petani"},
		{name: "Ferison", fig: 3, major: typE, minor: typI,
			stem: "Premis 1: Tidak ada dokter yang merupakan petani. Premis 2: Sebagian dokter adalah guru. Kesimpulan yang tepat adalah ...",
			key:  "Sebagian guru bukan petani"},
		{name: "Camenes", fig: 4, major: typA, minor: typE,
			stem: "Premis 1: Semua petani adalah dokter. Premis 2: Tidak ada dokter yang merupakan guru. Kesimpulan yang tepat adalah ...",
			key:  "Tidak ada guru yang merupakan petani"},
		{name: "Barbara dengan arah kesimpulan tertukar", fig: 4, major: typA, minor: typA,
			stem: "Premis 1: Semua petani adalah dokter. Premis 2: Semua dokter adalah guru. Kesimpulan yang tepat adalah ...",
			key:  "Semua petani adalah guru"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := params{kind: kindCat, names: [3]string{"guru", "dokter", "petani"}, prem: shape{c.fig, c.major, c.minor}.premises(), keyIdx: 0}
			res, ok := solve(p)
			if !ok {
				t.Fatalf("fixture params are rejected")
			}
			d := build(p, 7)
			if d.Stem != c.stem {
				t.Errorf("stem = %q", d.Stem)
			}
			var keyText string
			for _, o := range d.Options {
				if o.Key == d.AnswerKey {
					keyText = o.Text
				}
			}
			if keyText != c.key {
				t.Errorf("key = %q, want %q", keyText, c.key)
			}
			got := candTexts(res)
			for text, mis := range c.wantCands {
				if got[text] != mis {
					t.Errorf("candidate %q has mis %q, want %q", text, got[text], mis)
				}
			}
			for _, text := range c.absent {
				if _, found := got[text]; found {
					t.Errorf("candidate %q must not be offered", text)
				}
				for _, o := range d.Options {
					if o.Text == text {
						t.Errorf("option %q must not be offered", text)
					}
				}
			}
			if !strings.Contains(strings.Join(d.Steps, "\n"), c.key) {
				t.Errorf("steps %q lack the key", d.Steps)
			}
			band := gen.Band2
			if !strings.Contains(c.stem, "Sebagian") {
				band = gen.Band1
			}
			checkCategorical(t, d, band, 7)
		})
	}
}

func TestConditionalFixtures(t *testing.T) {
	// scenario 0: ids 0, 1, 2 = hujan turun, jalanan basah, jemuran diangkat
	a, b, c := lit{0, false}, lit{1, false}, lit{2, false}
	base := params{kind: kindCond, scenario: 0, ids: [3]int{0, 1, 2}}
	cases := []struct {
		name      string
		sch       scheme
		prem      [2]prop
		key       prop
		stem      string
		keyText   string
		wantCands map[string]string
		absent    []string
		stepPart  string
	}{
		{
			name: "modus ponens", sch: schemeMP, prem: [2]prop{imp(a, b), fact(a)}, key: fact(b),
			stem:    "Premis 1: Jika hujan turun maka jalanan basah. Premis 2: Hujan turun. Kesimpulan yang tepat adalah ...",
			keyText: "Jalanan basah",
			wantCands: map[string]string{
				"Jalanan tidak basah":                          MisNegasiKeliru,
				"Hujan tidak turun":                            MisNegasiKeliru,
				"Jika jalanan basah maka jemuran diangkat":     MisKesimpulanTakTerkait,
				"Jika hujan turun maka jemuran tidak diangkat": MisKesimpulanTakTerkait,
			},
			stepPart: "(modus ponens)",
		},
		{
			name: "modus tollens", sch: schemeMT, prem: [2]prop{fact(b.not()), imp(a, b)}, key: fact(a.not()),
			stem:    "Premis 1: Jalanan tidak basah. Premis 2: Jika hujan turun maka jalanan basah. Kesimpulan yang tepat adalah ...",
			keyText: "Hujan tidak turun",
			wantCands: map[string]string{
				"Hujan turun":   MisNegasiKeliru,
				"Jalanan basah": MisNegasiKeliru,
				"Jika hujan tidak turun maka jemuran diangkat":         MisKesimpulanTakTerkait,
				"Jika jalanan tidak basah maka jemuran tidak diangkat": MisKesimpulanTakTerkait,
			},
			stepPart: "(modus tollens)",
		},
		{
			name: "silogisme hipotetis", sch: schemeHS, prem: [2]prop{imp(a, b), imp(b, c)}, key: imp(a, c),
			stem:    "Premis 1: Jika hujan turun maka jalanan basah. Premis 2: Jika jalanan basah maka jemuran diangkat. Kesimpulan yang tepat adalah ...",
			keyText: "Jika hujan turun maka jemuran diangkat",
			wantCands: map[string]string{
				"Jika jemuran diangkat maka hujan turun":             MisAfirmasiKonsekuen,
				"Jika hujan tidak turun maka jemuran tidak diangkat": MisSangkalAnteseden,
				"Jika jalanan basah maka hujan turun":                MisKonversiKeliru,
				"Jika jemuran diangkat maka jalanan basah":           MisKonversiKeliru,
				"Jika hujan turun maka jemuran tidak diangkat":       MisNegasiKeliru,
				"Jika hujan tidak turun maka jemuran diangkat":       MisNegasiKeliru,
				"Hujan turun":            MisKesimpulanTakTerkait,
				"Jemuran tidak diangkat": MisKesimpulanTakTerkait,
			},
			// the contrapositive of the key is also valid, so it must not be offered
			absent:   []string{"Jika jemuran tidak diangkat maka hujan tidak turun"},
			stepPart: "silogisme hipotetis",
		},
		{
			name: "silogisme hipotetis, bentuk kontraposisi", sch: schemeHS, prem: [2]prop{contra(imp(a, b)), imp(b, c)}, key: contra(imp(a, c)),
			stem:    "Premis 1: Jika jalanan tidak basah maka hujan tidak turun. Premis 2: Jika jalanan basah maka jemuran diangkat. Kesimpulan yang tepat adalah ...",
			keyText: "Jika jemuran tidak diangkat maka hujan tidak turun",
			wantCands: map[string]string{
				"Jika hujan tidak turun maka jemuran tidak diangkat": MisAfirmasiKonsekuen,
				"Jika jemuran diangkat maka hujan turun":             MisSangkalAnteseden,
			},
			// the plain form of the key is also valid, so it must not be offered
			absent:   []string{"Jika hujan turun maka jemuran diangkat"},
			stepPart: "kontraposisinya",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := base
			p.scheme, p.cprem, p.ckey = tc.sch, tc.prem, tc.key
			res, ok := solve(p)
			if !ok {
				t.Fatalf("fixture params are rejected")
			}
			d := build(p, 7)
			if d.Stem != tc.stem {
				t.Errorf("stem = %q", d.Stem)
			}
			var keyText string
			for _, o := range d.Options {
				if o.Key == d.AnswerKey {
					keyText = o.Text
				}
			}
			if keyText != tc.keyText {
				t.Errorf("key = %q, want %q", keyText, tc.keyText)
			}
			got := candTexts(res)
			for text, mis := range tc.wantCands {
				if got[text] != mis {
					t.Errorf("candidate %q has mis %q, want %q", text, got[text], mis)
				}
			}
			for _, text := range tc.absent {
				if _, found := got[text]; found {
					t.Errorf("candidate %q must not be offered", text)
				}
			}
			if !strings.Contains(strings.Join(d.Steps, "\n"), tc.stepPart) {
				t.Errorf("steps %q lack %q", d.Steps, tc.stepPart)
			}
			checkConditional(t, d, 7)
		})
	}
}

// ---- variants, determinism, errors ------------------------------------------

func TestBandsEmitExpectedVariants(t *testing.T) {
	// band 3: one premise is a fact (modus ponens or tollens) or both are
	// conditionals (hypothetical syllogism); the key is a fact or a conditional
	seen := map[string]int{}
	for seed := int64(1); seed <= 500; seed++ {
		d, err := Family{}.Generate(seed, gen.Band3)
		if err != nil {
			t.Fatal(err)
		}
		s1, s2 := splitStem(t, d.Stem)
		conds := strings.Count(s1+"|"+s2, "Jika ")
		seen[fmt.Sprint("premises with Jika: ", conds)]++
		for _, o := range d.Options {
			if o.Key == d.AnswerKey {
				seen[fmt.Sprint("key conditional: ", strings.HasPrefix(o.Text, "Jika "))]++
			}
		}
	}
	for _, k := range []string{"premises with Jika: 1", "premises with Jika: 2", "key conditional: true", "key conditional: false"} {
		if seen[k] == 0 {
			t.Errorf("band 3 never emitted %q (seen %v)", k, seen)
		}
	}
	if len(seen) != 4 {
		t.Errorf("band 3 emitted unexpected variants %v", seen)
	}
}

func TestCategoricalBandsEmitExpectedVariants(t *testing.T) {
	premTypes := map[gen.Band]map[byte]bool{gen.Band1: {}, gen.Band2: {}}
	keyTypes := map[gen.Band]map[byte]bool{gen.Band1: {}, gen.Band2: {}}
	for band := gen.Band1; band <= gen.Band2; band++ {
		for seed := int64(1); seed <= 1000; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			s1, s2 := splitStem(t, d.Stem)
			for _, s := range []string{s1, s2} {
				p, _ := parseCat(s)
				premTypes[band][p.typ] = true
			}
			for _, o := range d.Options {
				if o.Key == d.AnswerKey {
					k, _ := parseCat(o.Text)
					keyTypes[band][k.typ] = true
				}
			}
		}
	}
	if got := premTypes[gen.Band1]; len(got) != 2 || !got['A'] || !got['E'] {
		t.Errorf("band 1 premise types = %v, want exactly A and E", got)
	}
	if got := premTypes[gen.Band2]; !got['I'] || !got['O'] {
		t.Errorf("band 2 premise types = %v, want I and O among them", got)
	}
	if got := keyTypes[gen.Band1]; !got['A'] || !got['E'] {
		t.Errorf("band 1 key types = %v, want A and E", got)
	}
	if got := keyTypes[gen.Band2]; !got['I'] || !got['O'] {
		t.Errorf("band 2 key types = %v, want I and O", got)
	}
}

func TestEveryMisconceptionIsEmitted(t *testing.T) {
	seen := map[string]bool{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 2000; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			for _, o := range d.Options {
				seen[o.MisconceptionID] = true
			}
		}
	}
	for id := range (Family{}).Misconceptions() {
		if !seen[id] {
			t.Errorf("misconception %q is never emitted", id)
		}
	}
}

func TestMisconceptionLabelsCoverEmittedIds(t *testing.T) {
	fam := Family{}
	labels := fam.Misconceptions()
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 500; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			for _, o := range d.Options {
				if o.MisconceptionID != "" && labels[o.MisconceptionID] == "" {
					t.Fatalf("band %d seed %d: emitted id %q has no label", band, seed, o.MisconceptionID)
				}
			}
		}
	}
	if fam.ID() != PatternID || fam.Label() != "Silogisme" || PatternID != "TIU.verbal.silogisme" {
		t.Errorf("id/label = %q/%q", fam.ID(), fam.Label())
	}
	labels["x"] = "y"
	if _, ok := fam.Misconceptions()["x"]; ok {
		t.Error("Misconceptions returns the shared map")
	}
}

func TestDeterminismAndVariation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		stems := map[string]bool{}
		for seed := int64(1); seed <= 200; seed++ {
			a, err1 := fam.Generate(seed, band)
			b, err2 := fam.Generate(seed, band)
			if err1 != nil || err2 != nil {
				t.Fatal(err1, err2)
			}
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("band %d seed %d: drafts differ", band, seed)
			}
			stems[a.Stem] = true
		}
		t.Logf("band %d: %d distinct stems in 200 seeds", band, len(stems))
		if len(stems) < 150 {
			t.Errorf("band %d: only %d distinct stems in 200 seeds", band, len(stems))
		}
	}
}

func TestSolveRejectsMalformedParams(t *testing.T) {
	names := [3]string{"guru", "dokter", "petani"}
	barbara := shape{1, typA, typA}.premises()
	a, b, c := lit{0, false}, lit{1, false}, lit{2, false}
	cases := map[string]params{
		"tanpa kunci sah (AA figur 3)":      {kind: kindCat, names: names, prem: shape{3, typA, typA}.premises()},
		"tanpa kunci sah (EE)":              {kind: kindCat, names: names, prem: shape{1, typE, typE}.premises()},
		"tanpa kunci sah (II)":              {kind: kindCat, names: names, prem: shape{1, typI, typI}.premises()},
		"premis sama istilah":               {kind: kindCat, names: names, prem: [2]stmt{{typA, 1, 1}, barbara[1]}},
		"istilah di luar 0..2":              {kind: kindCat, names: names, prem: [2]stmt{{typA, 1, 3}, barbara[1]}},
		"tipe tak dikenal":                  {kind: kindCat, names: names, prem: [2]stmt{{stmtType(9), 1, 2}, barbara[1]}},
		"nama kosong":                       {kind: kindCat, names: [3]string{"guru", "", "petani"}, prem: barbara},
		"keyIdx negatif":                    {kind: kindCat, names: names, prem: barbara, keyIdx: -1},
		"premis kosong secara logis":        {kind: kindCat, names: names, prem: [2]stmt{{typI, 1, 2}, {typE, 1, 2}}},
		"kondisional kunci tidak sah":       {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{imp(a, b), fact(b)}, ckey: fact(a)},
		"kondisional kunci tak terkait":     {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{imp(a, b), fact(a)}, ckey: fact(c)},
		"kondisional premis kontradiktif":   {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{fact(a), fact(a.not())}, ckey: fact(b)},
		"kondisional id sama":               {kind: kindCond, ids: [3]int{0, 0, 2}, scheme: schemeMP, cprem: [2]prop{imp(a, b), fact(a)}, ckey: fact(b)},
		"kondisional id di luar kolam":      {kind: kindCond, ids: [3]int{0, 1, 9}, scheme: schemeMP, cprem: [2]prop{imp(a, b), fact(a)}, ckey: fact(b)},
		"kondisional skenario tak dikenal":  {kind: kindCond, scenario: 99, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{imp(a, b), fact(a)}, ckey: fact(b)},
		"kondisional rantai tanpa kondisi":  {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeHS, cprem: [2]prop{imp(a, b), fact(b)}, ckey: imp(a, b)},
		"kondisional skema tak dikenal":     {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: scheme(9), cprem: [2]prop{imp(a, b), fact(a)}, ckey: fact(b)},
		"kondisional dua fakta":             {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{fact(a), fact(a)}, ckey: fact(a)},
		"kondisional variabel di luar 0..2": {kind: kindCond, ids: [3]int{0, 1, 2}, scheme: schemeMP, cprem: [2]prop{imp(a, lit{5, false}), fact(a)}, ckey: fact(lit{5, false})},
		"jenis tidak dikenal":               {kind: kind(9)},
	}
	for name, p := range cases {
		if _, ok := solve(p); ok {
			t.Errorf("%s: solve accepted %+v", name, p)
		}
	}
}

func TestGenerateErrorsWithoutValidParams(t *testing.T) {
	if _, err := generate(1, gen.Band2, maxAttempts, func(params) bool { return false }); err == nil {
		t.Fatal("expected error when no params are accepted")
	}
	for _, band := range []gen.Band{0, 4, 9} {
		if _, err := (Family{}).Generate(1, band); err == nil {
			t.Fatalf("expected error for unknown band %d", band)
		}
	}
}

func TestDistinctStemsPer20000(t *testing.T) {
	// measured: 13493 / 17140 / 12591
	floor := map[gen.Band]int{gen.Band1: 13400, gen.Band2: 17000, gen.Band3: 12500}
	for band := gen.Band1; band <= gen.Band3; band++ {
		stems := map[string]bool{}
		for seed := int64(1); seed <= 20000; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			stems[d.Stem] = true
		}
		t.Logf("band %d: %d distinct stems in 20000 seeds", band, len(stems))
		if len(stems) < floor[band] {
			t.Errorf("band %d: only %d distinct stems in 20000 seeds, want >= %d", band, len(stems), floor[band])
		}
	}
}
