// Package silogisme generates syllogism questions (pattern
// TIU.verbal.silogisme). Bands 1 and 2 are categorical syllogisms, band 3 is
// conditional ("Jika ... maka ...") reasoning.
//
// The key is never taken from a table of forms. Validity is decided by
// enumerating finite models:
//
//   - categorical: three terms S, M, P split a universe into 8 Venn regions; a
//     model is the set of non-empty regions (2^8 models). A conclusion is valid
//     when it is true in every model that satisfies both premises.
//   - conditional: three propositions give 8 truth assignments; a conclusion is
//     valid when it is true in every assignment that satisfies both premises
//     (material implication).
//
// Convention: "Semua X adalah Y" does not guarantee that X has members, and
// "Sebagian" means at least one. The key must be valid under that reading, and
// every wrong option must fail even when all three terms are non-empty and even
// when "Sebagian X adalah Y" is read as "sebagian saja" (some X are also not Y,
// and the other way round), so no option depends on whether "semua" implies
// "sebagian" or on how "sebagian" is read.
package silogisme

import (
	"fmt"
	"math/rand"
	"strings"

	"latihseleksi/api/internal/gen"
)

const (
	PatternID   = "TIU.verbal.silogisme"
	maxAttempts = 50
	shuffleSalt = 1 << 20
	optionCount = 5
)

// Misconception labels attached to wrong options.
const (
	MisKonversiKeliru       = "konversi-keliru"
	MisSebagianJadiSemua    = "sebagian-jadi-semua"
	MisNegasiKeliru         = "negasi-keliru"
	MisAfirmasiKonsekuen    = "afirmasi-konsekuen"
	MisSangkalAnteseden     = "sangkal-anteseden"
	MisKesimpulanTakTerkait = "kesimpulan-tak-terkait"
)

var misconceptionLabels = map[string]string{
	MisKonversiKeliru:       "Membalik arah hubungan begitu saja: yang berlaku dari X ke Y dianggap berlaku juga dari Y ke X (pada \"Semua ... adalah ...\" maupun \"Jika ... maka ...\")",
	MisSebagianJadiSemua:    "Menyimpulkan lebih kuat dari premis: \"sebagian\" dianggap \"semua\", atau \"ada yang bukan\" dianggap \"tidak ada\"",
	MisNegasiKeliru:         "Memilih kesimpulan yang berlawanan dengan hubungan yang sebenarnya (positif menjadi negatif atau sebaliknya)",
	MisAfirmasiKonsekuen:    "Menganggap syarat (anteseden) pasti benar karena akibatnya (konsekuen) benar, yaitu membalik arah \"jika-maka\"",
	MisSangkalAnteseden:     "Menganggap akibat pasti salah karena syaratnya salah",
	MisKesimpulanTakTerkait: "Memilih kalimat yang tidak mengikuti dari kedua premis: istilah penghubung ikut disimpulkan, fakta tunggal, atau hal yang tidak ditentukan premis",
}

// kind selects the question variant.
type kind int

const (
	kindCat  kind = iota // categorical syllogism (bands 1 and 2)
	kindCond             // conditional reasoning (band 3)
)

// ---- categorical syllogisms -------------------------------------------------

// stmtType is the traditional type of a categorical statement.
type stmtType int

const (
	typA stmtType = iota // Semua x adalah y
	typE                 // Tidak ada x yang merupakan y
	typI                 // Sebagian x adalah y
	typO                 // Sebagian x bukan y
)

func (t stmtType) negative() bool { return t == typE || t == typO }

// Term roles: S and P are the ends of the conclusion, M is the middle term.
const (
	roleS = 0
	roleM = 1
	roleP = 2
)

// stmt relates two terms (role ids).
type stmt struct {
	t    stmtType
	x, y int
}

// holds reports whether s is true in a model: bit r of model is set when
// Venn region r is non-empty, and bit k of r says the region lies inside term k.
func (s stmt) holds(model uint8) bool { return truth[s.t][s.x][s.y][model] }

// truth caches stmt.holds for every statement and model.
var truth = buildTruth()

func buildTruth() (tab [4][3][3][256]bool) {
	for t := typA; t <= typO; t++ {
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				for m := 0; m < 256; m++ {
					tab[t][x][y][m] = evalStmt(stmt{t, x, y}, uint8(m))
				}
			}
		}
	}
	return tab
}

func evalStmt(s stmt, model uint8) bool {
	for r := 0; r < 8; r++ {
		if model>>r&1 == 0 {
			continue
		}
		inX, inY := r>>s.x&1 == 1, r>>s.y&1 == 1
		switch s.t {
		case typA:
			if inX && !inY {
				return false
			}
		case typE:
			if inX && inY {
				return false
			}
		case typI:
			if inX && inY {
				return true
			}
		case typO:
			if inX && !inY {
				return true
			}
		}
	}
	return s.t == typA || s.t == typE
}

// termsNonEmpty reports whether every term has a member in model.
func termsNonEmpty(model uint8) bool {
	for k := 0; k < 3; k++ {
		found := false
		for r := 0; r < 8; r++ {
			if model>>r&1 == 1 && r>>k&1 == 1 {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func allHold(prem []stmt, model uint8) bool {
	for _, s := range prem {
		if !s.holds(model) {
			return false
		}
	}
	return true
}

// premiseModels returns the models that satisfy every premise.
func premiseModels(prem []stmt) []uint8 {
	var out []uint8
	for m := 0; m < 256; m++ {
		if allHold(prem, uint8(m)) {
			out = append(out, uint8(m))
		}
	}
	return out
}

// onlyNonEmpty keeps the models in which all three terms have members.
func onlyNonEmpty(models []uint8) []uint8 {
	var out []uint8
	for _, m := range models {
		if termsNonEmpty(m) {
			out = append(out, m)
		}
	}
	return out
}

// validIn reports whether c is true in every model; with no model the premises
// prove nothing, so false.
func validIn(models []uint8, c stmt) bool {
	for _, m := range models {
		if !c.holds(m) {
			return false
		}
	}
	return len(models) > 0
}

// refutedIn reports whether c is false in some model.
func refutedIn(models []uint8, c stmt) bool {
	for _, m := range models {
		if !c.holds(m) {
			return true
		}
	}
	return false
}

// catValid reports whether c is true in every model of the premises (no
// existential import).
func catValid(prem []stmt, c stmt) bool { return validIn(premiseModels(prem), c) }

// catRefuted reports whether some model with all three terms non-empty
// satisfies the premises and falsifies c.
func catRefuted(prem []stmt, c stmt) bool {
	return refutedIn(onlyNonEmpty(premiseModels(prem)), c)
}

// catConsistent reports whether the premises have a model with all terms non-empty.
func catConsistent(prem []stmt) bool { return len(onlyNonEmpty(premiseModels(prem))) > 0 }

// implicature adds the reading in which "Sebagian X adalah Y" also says that
// some X are not Y, and the other way round.
func implicature(prem []stmt) []stmt {
	out := append([]stmt(nil), prem...)
	for _, s := range prem {
		switch s.t {
		case typI:
			out = append(out, stmt{typO, s.x, s.y})
		case typO:
			out = append(out, stmt{typI, s.x, s.y})
		}
	}
	return out
}

// shape is a pair of premise types in one of the four figures.
type shape struct {
	fig          int // 1..4: where M stands in the major and minor premise
	major, minor stmtType
}

// premises returns the major premise (with P) and the minor premise (with S).
func (s shape) premises() [2]stmt {
	major, minor := stmt{s.major, roleP, roleM}, stmt{s.minor, roleM, roleS}
	if s.fig == 1 || s.fig == 3 {
		major = stmt{s.major, roleM, roleP}
	}
	if s.fig == 1 || s.fig == 2 {
		minor = stmt{s.minor, roleS, roleM}
	}
	return [2]stmt{major, minor}
}

// endConclusions lists the 8 statements relating S and P.
func endConclusions() []stmt {
	var out []stmt
	for t := typA; t <= typO; t++ {
		out = append(out, stmt{t, roleS, roleP}, stmt{t, roleP, roleS})
	}
	return out
}

// middleConclusions lists the 16 statements relating M to S or P.
func middleConclusions() []stmt {
	var out []stmt
	for t := typA; t <= typO; t++ {
		for _, o := range []int{roleS, roleP} {
			out = append(out, stmt{t, o, roleM}, stmt{t, roleM, o})
		}
	}
	return out
}

func validEndConclusions(prem [2]stmt) []stmt {
	var out []stmt
	for _, c := range endConclusions() {
		if catValid(prem[:], c) {
			out = append(out, c)
		}
	}
	return out
}

// catShapes lists, per band, the premise shapes that have a valid conclusion
// between S and P: band 1 has only Semua/Tidak ada premises, band 2 at least
// one Sebagian premise.
var catShapes = buildShapes()

func buildShapes() map[gen.Band][]shape {
	out := map[gen.Band][]shape{}
	for fig := 1; fig <= 4; fig++ {
		for a := typA; a <= typO; a++ {
			for b := typA; b <= typO; b++ {
				sh := shape{fig, a, b}
				prem := sh.premises()
				if len(validEndConclusions(prem)) == 0 || !catConsistent(prem[:]) {
					continue
				}
				band := gen.Band2
				if a <= typE && b <= typE {
					band = gen.Band1
				}
				out[band] = append(out[band], sh)
			}
		}
	}
	return out
}

// termGroups hold neutral terms; the three terms of a question share a group
// so that real-world knowledge does not decide the answer.
var termGroups = [][]string{
	{"guru", "dokter", "petani", "nelayan", "pedagang", "sopir", "montir", "penjahit", "tukang kayu", "arsitek", "apoteker", "pilot"},
	{"pelukis", "pemancing", "pendaki", "pelari", "perenang", "pembaca", "penari", "pemain catur", "pemain gitar", "fotografer", "pesepeda", "kolektor perangko"},
	{"relawan", "peserta pelatihan", "anggota koperasi", "peserta lomba", "pengunjung perpustakaan", "penghuni asrama", "anggota paduan suara", "panitia acara", "penonton pertunjukan", "peserta seminar", "pelanggan toko", "anggota klub"},
}

// ---- conditional reasoning --------------------------------------------------

type proposition struct{ pos, neg string }

// scenarios hold propositions that may be combined within one scenario.
var scenarios = [][]proposition{
	{
		{"hujan turun", "hujan tidak turun"},
		{"jalanan basah", "jalanan tidak basah"},
		{"jemuran diangkat", "jemuran tidak diangkat"},
		{"sungai meluap", "sungai tidak meluap"},
		{"sekolah diliburkan", "sekolah tidak diliburkan"},
		{"lapangan tergenang", "lapangan tidak tergenang"},
		{"pengendara memakai jas hujan", "pengendara tidak memakai jas hujan"},
	},
	{
		{"lampu ruang rapat menyala", "lampu ruang rapat tidak menyala"},
		{"proyektor aktif", "proyektor tidak aktif"},
		{"rapat dimulai", "rapat tidak dimulai"},
		{"notulen dibuat", "notulen tidak dibuat"},
		{"pintu ruang rapat terkunci", "pintu ruang rapat tidak terkunci"},
		{"daftar hadir diisi", "daftar hadir tidak diisi"},
		{"pendingin ruangan dinyalakan", "pendingin ruangan tidak dinyalakan"},
	},
	{
		{"siswa mengerjakan latihan", "siswa tidak mengerjakan latihan"},
		{"siswa memahami materi", "siswa tidak memahami materi"},
		{"siswa lulus ujian", "siswa tidak lulus ujian"},
		{"siswa mendapat beasiswa", "siswa tidak mendapat beasiswa"},
		{"guru memberi pujian", "guru tidak memberi pujian"},
		{"perpustakaan ramai", "perpustakaan tidak ramai"},
		{"siswa mengulang pelajaran", "siswa tidak mengulang pelajaran"},
	},
	{
		{"mesin dihidupkan", "mesin tidak dihidupkan"},
		{"aki terisi penuh", "aki tidak terisi penuh"},
		{"lampu indikator menyala", "lampu indikator tidak menyala"},
		{"bensin tersedia", "bensin tidak tersedia"},
		{"mobil dapat berjalan", "mobil tidak dapat berjalan"},
		{"rem diperiksa", "rem tidak diperiksa"},
		{"bengkel buka", "bengkel tidak buka"},
	},
	{
		{"pelanggan datang", "pelanggan tidak datang"},
		{"stok barang tersedia", "stok barang tidak tersedia"},
		{"kasir bertugas", "kasir tidak bertugas"},
		{"promo diadakan", "promo tidak diadakan"},
		{"toko tutup lebih awal", "toko tidak tutup lebih awal"},
		{"pesanan dikirim", "pesanan tidak dikirim"},
		{"pembayaran diterima", "pembayaran tidak diterima"},
	},
	{
		{"taman dibersihkan", "taman tidak dibersihkan"},
		{"relawan hadir", "relawan tidak hadir"},
		{"acara berlangsung", "acara tidak berlangsung"},
		{"panitia menyiapkan tenda", "panitia tidak menyiapkan tenda"},
		{"peserta mendaftar", "peserta tidak mendaftar"},
		{"cuaca cerah", "cuaca tidak cerah"},
		{"pengunjung membeli tiket", "pengunjung tidak membeli tiket"},
	},
}

// lit is proposition v (0..2) or its negation.
type lit struct {
	v   int
	neg bool
}

func (l lit) not() lit { return lit{l.v, !l.neg} }

// prop is a fact (cond false, only a) or "Jika a maka b".
type prop struct {
	cond bool
	a, b lit
}

func fact(l lit) prop   { return prop{a: l} }
func imp(a, b lit) prop { return prop{cond: true, a: a, b: b} }

// contra is the contrapositive of a conditional.
func contra(q prop) prop { return imp(q.b.not(), q.a.not()) }

// holds evaluates p under assignment w: bit v of w is the truth value of proposition v.
func (p prop) holds(w int) bool {
	val := func(l lit) bool { return (w>>l.v&1 == 1) != l.neg }
	if !p.cond {
		return val(p.a)
	}
	return !val(p.a) || val(p.b)
}

func propsHold(ps []prop, w int) bool {
	for _, p := range ps {
		if !p.holds(w) {
			return false
		}
	}
	return true
}

// condValid reports whether c is true in every assignment that satisfies the
// premises; premises without any such assignment prove nothing, so false.
func condValid(prem []prop, c prop) bool {
	sat := false
	for w := 0; w < 8; w++ {
		if !propsHold(prem, w) {
			continue
		}
		sat = true
		if !c.holds(w) {
			return false
		}
	}
	return sat
}

// condRefuted reports whether some assignment satisfies the premises and falsifies c.
func condRefuted(prem []prop, c prop) bool {
	for w := 0; w < 8; w++ {
		if propsHold(prem, w) && !c.holds(w) {
			return true
		}
	}
	return false
}

func condConsistent(prem []prop) bool {
	for w := 0; w < 8; w++ {
		if propsHold(prem, w) {
			return true
		}
	}
	return false
}

// scheme is the inference rule behind a conditional question.
type scheme int

const (
	schemeMP scheme = iota // fact a, premise a -> b: b
	schemeMT               // fact not b, premise a -> b: not a
	schemeHS               // a -> b and b -> c: a -> c
)

// ---- shared -----------------------------------------------------------------

// params are the choices behind one question.
//
//	kindCat: names holds the terms by role (S, M, P), prem the premises in
//	  display order and keyIdx picks one of the valid S-P conclusions.
//	kindCond: scenario and ids pick the three propositions, prem holds the
//	  premises in display order and ckey the intended conclusion.
type params struct {
	kind     kind
	names    [3]string
	prem     [2]stmt
	keyIdx   int
	scenario int
	ids      [3]int
	cprem    [2]prop
	ckey     prop
	scheme   scheme
}

// Family implements gen.Family for syllogism questions.
type Family struct{}

var _ gen.Family = Family{}

func (Family) ID() string { return PatternID }

func (Family) Label() string { return "Silogisme" }

// Misconceptions returns the labels of the ids in the Mis* constants.
func (Family) Misconceptions() map[string]string {
	out := make(map[string]string, len(misconceptionLabels))
	for k, v := range misconceptionLabels {
		out[k] = v
	}
	return out
}

// Generate returns the question for seed and band. Parameters that do not
// give a valid key and enough wrong options are rejected; attempt k uses a
// seed derived from (seed, k), and after maxAttempts it returns an error.
func (Family) Generate(seed int64, band gen.Band) (gen.Draft, error) {
	return generate(seed, band, maxAttempts, valid)
}

func valid(p params) bool {
	_, ok := solve(p)
	return ok
}

func generate(seed int64, band gen.Band, attempts int, accept func(params) bool) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("silogisme: band %d tidak dikenal", band)
	}
	for k := 0; k < attempts; k++ {
		rng := rand.New(rand.NewSource(derive(seed, k)))
		p := pickParams(rng, band)
		if accept(p) {
			return build(p, seed), nil
		}
	}
	return gen.Draft{}, fmt.Errorf("silogisme: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", seed, band, attempts)
}

// derive mixes (seed, k) into a new seed (splitmix64).
func derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

func pickParams(rng *rand.Rand, band gen.Band) params {
	if band == gen.Band3 {
		return pickCond(rng)
	}
	return pickCat(rng, band)
}

func pickCat(rng *rand.Rand, band gen.Band) params {
	shapes := catShapes[band]
	sh := shapes[rng.Intn(len(shapes))]
	g := termGroups[rng.Intn(len(termGroups))]
	perm := rng.Perm(len(g))
	prem := sh.premises()
	if rng.Intn(2) == 0 {
		prem[0], prem[1] = prem[1], prem[0]
	}
	return params{
		kind:   kindCat,
		names:  [3]string{g[perm[0]], g[perm[1]], g[perm[2]]},
		prem:   prem,
		keyIdx: rng.Intn(8),
	}
}

func pickCond(rng *rand.Rand) params {
	sc := rng.Intn(len(scenarios))
	perm := rng.Perm(len(scenarios[sc]))
	a, b, c := lit{0, rng.Intn(2) == 0}, lit{1, rng.Intn(2) == 0}, lit{2, rng.Intn(2) == 0}
	// form writes a conditional directly or as its contrapositive
	form := func(x, y lit) prop {
		if rng.Intn(2) == 0 {
			return imp(x, y)
		}
		return contra(imp(x, y))
	}
	p := params{kind: kindCond, scenario: sc, ids: [3]int{perm[0], perm[1], perm[2]}}
	switch rng.Intn(4) {
	case 0:
		p.scheme, p.cprem, p.ckey = schemeMP, [2]prop{form(a, b), fact(a)}, fact(b)
	case 1:
		p.scheme, p.cprem, p.ckey = schemeMT, [2]prop{form(a, b), fact(b.not())}, fact(a.not())
	default:
		p.scheme, p.cprem = schemeHS, [2]prop{form(a, b), form(b, c)}
		p.ckey = imp(a, c)
		if rng.Intn(2) == 0 {
			p.ckey = contra(p.ckey)
		}
	}
	if rng.Intn(2) == 0 {
		p.cprem[0], p.cprem[1] = p.cprem[1], p.cprem[0]
	}
	return p
}

type cand struct {
	text string
	mis  string
}

// result is a solved question: the key text, the wrong options that were
// proven invalid, and the solution steps.
type result struct {
	key   string
	cands []cand
	steps []string
}

// solve decides the key and collects wrong options, each proven invalid by the
// model checker. ok is false when the parameters are malformed, the premises
// are inconsistent, the key is not valid, or fewer than optionCount-1 wrong
// options exist.
func solve(p params) (result, bool) {
	var res result
	var ok bool
	switch p.kind {
	case kindCat:
		res, ok = solveCat(p)
	case kindCond:
		res, ok = solveCond(p)
	}
	if !ok {
		return result{}, false
	}
	seen := map[string]bool{res.key: true}
	var uniq []cand
	for _, c := range res.cands {
		if !seen[c.text] {
			seen[c.text] = true
			uniq = append(uniq, c)
		}
	}
	res.cands = uniq
	return res, len(uniq) >= optionCount-1
}

func (p params) stmtText(s stmt) string {
	x, y := p.names[s.x], p.names[s.y]
	switch s.t {
	case typA:
		return fmt.Sprintf("Semua %s adalah %s", x, y)
	case typE:
		return fmt.Sprintf("Tidak ada %s yang merupakan %s", x, y)
	case typI:
		return fmt.Sprintf("Sebagian %s adalah %s", x, y)
	}
	return fmt.Sprintf("Sebagian %s bukan %s", x, y)
}

func (p params) vennText(s stmt) string {
	x, y := p.names[s.x], p.names[s.y]
	switch s.t {
	case typA:
		return fmt.Sprintf("seluruh lingkaran %s berada di dalam lingkaran %s", x, y)
	case typE:
		return fmt.Sprintf("lingkaran %s dan %s tidak beririsan", x, y)
	case typI:
		return fmt.Sprintf("lingkaran %s dan %s beririsan dan irisannya berisi anggota", x, y)
	}
	return fmt.Sprintf("ada anggota %s di luar lingkaran %s", x, y)
}

// classify names the mistake behind a wrong end conclusion c relative to the key.
func classify(key, c stmt) string {
	switch {
	case c.t == key.t && c.x == key.y && c.y == key.x:
		return MisKonversiKeliru
	case c.t.negative() != key.t.negative():
		return MisNegasiKeliru
	case (key.t == typI && c.t == typA) || (key.t == typO && c.t == typE):
		return MisSebagianJadiSemua
	}
	return MisKesimpulanTakTerkait
}

func solveCat(p params) (result, bool) {
	for _, s := range p.prem {
		if s.t < typA || s.t > typO || s.x < 0 || s.x > 2 || s.y < 0 || s.y > 2 || s.x == s.y {
			return result{}, false
		}
	}
	for _, n := range p.names {
		if n == "" {
			return result{}, false
		}
	}
	prem := p.prem[:]
	if p.keyIdx < 0 {
		return result{}, false
	}
	// A wrong option must fail with all terms non-empty, and also when
	// "Sebagian" is read as "sebagian saja"; the key must hold in every model.
	plain := onlyNonEmpty(premiseModels(prem))
	strong := onlyNonEmpty(premiseModels(implicature(prem)))
	if len(plain) == 0 || len(strong) == 0 {
		return result{}, false
	}
	all := premiseModels(prem)
	var valids []stmt
	for _, c := range endConclusions() {
		if validIn(all, c) {
			valids = append(valids, c)
		}
	}
	if len(valids) == 0 {
		return result{}, false
	}
	key := valids[p.keyIdx%len(valids)]
	res := result{key: p.stmtText(key)}
	for _, c := range endConclusions() {
		if c != key && refutedIn(plain, c) && refutedIn(strong, c) {
			res.cands = append(res.cands, cand{p.stmtText(c), classify(key, c)})
		}
	}
	for _, c := range middleConclusions() {
		if refutedIn(plain, c) && refutedIn(strong, c) {
			res.cands = append(res.cands, cand{p.stmtText(c), MisKesimpulanTakTerkait})
		}
	}
	res.steps = []string{
		fmt.Sprintf("Premis 1 (%s): %s.", p.stmtText(p.prem[0]), p.vennText(p.prem[0])),
		fmt.Sprintf("Premis 2 (%s): %s.", p.stmtText(p.prem[1]), p.vennText(p.prem[1])),
		fmt.Sprintf("Kesimpulan sah hanya jika benar pada setiap gambar yang memenuhi kedua premis. Antara %s dan %s, yang selalu benar adalah: %s.", p.names[roleS], p.names[roleP], res.key),
		"Pilihan lain dapat digambar salah tanpa melanggar kedua premis.",
	}
	return res, true
}

func (p params) litText(l lit) string {
	pr := scenarios[p.scenario][p.ids[l.v]]
	if l.neg {
		return pr.neg
	}
	return pr.pos
}

// propText is the sentence for q; facts start with a capital letter.
func (p params) propText(q prop) string {
	if q.cond {
		return "Jika " + p.litText(q.a) + " maka " + p.litText(q.b)
	}
	return capitalize(p.litText(q.a))
}

func solveCond(p params) (result, bool) {
	if p.scenario < 0 || p.scenario >= len(scenarios) {
		return result{}, false
	}
	n := len(scenarios[p.scenario])
	if p.ids[0] == p.ids[1] || p.ids[0] == p.ids[2] || p.ids[1] == p.ids[2] {
		return result{}, false
	}
	for _, id := range p.ids {
		if id < 0 || id >= n {
			return result{}, false
		}
	}
	prem := p.cprem[:]
	for _, q := range append([]prop{p.ckey}, prem...) {
		if q.a.v < 0 || q.a.v > 2 || (q.cond && (q.b.v < 0 || q.b.v > 2 || q.a.v == q.b.v)) {
			return result{}, false
		}
	}
	if !condConsistent(prem) || !condValid(prem, p.ckey) {
		return result{}, false
	}
	res := result{key: p.propText(p.ckey)}
	add := func(c prop, mis string) {
		if condRefuted(prem, c) {
			res.cands = append(res.cands, cand{p.propText(c), mis})
		}
	}
	switch p.scheme {
	case schemeMP, schemeMT:
		var f prop
		var q prop
		found := 0
		for _, s := range p.cprem {
			if s.cond {
				q = s
			} else {
				f = s
				found++
			}
		}
		if found != 1 || p.ckey.cond || f.a.v == p.ckey.a.v {
			return result{}, false
		}
		third := lit{3 - f.a.v - p.ckey.a.v, false}
		add(fact(p.ckey.a.not()), MisNegasiKeliru)
		add(fact(f.a.not()), MisNegasiKeliru)
		// options that tie the third proposition to the fact or the key, so every
		// option talks about the question's own propositions
		add(imp(p.ckey.a, third), MisKesimpulanTakTerkait)
		add(imp(third, p.ckey.a), MisKesimpulanTakTerkait)
		add(imp(f.a, third.not()), MisKesimpulanTakTerkait)
		add(imp(third.not(), f.a), MisKesimpulanTakTerkait)
		// the conditional, as written or as its contrapositive, whose antecedent is the fact
		use, rule := q, "modus ponens"
		if use.a != f.a {
			use, rule = contra(q), "modus tollens"
		}
		if use.a != f.a || use.b != p.ckey.a {
			return result{}, false
		}
		res.steps = []string{
			fmt.Sprintf("Premis jika-maka \"%s\" setara dengan kontraposisinya \"%s\".", p.propText(q), p.propText(contra(q))),
			fmt.Sprintf("Fakta \"%s\" memenuhi syarat \"%s\", sehingga \"%s\" pasti benar (%s).", p.propText(f), p.propText(use), p.propText(p.ckey), rule),
		}
	case schemeHS:
		if !p.ckey.cond {
			return result{}, false
		}
		x, y := p.ckey.a, p.ckey.b
		add(imp(y, x), MisAfirmasiKonsekuen)
		add(imp(x.not(), y.not()), MisSangkalAnteseden)
		for _, q := range p.cprem {
			if !q.cond {
				return result{}, false
			}
			add(imp(q.b, q.a), MisKonversiKeliru)
		}
		add(imp(x, y.not()), MisNegasiKeliru)
		add(imp(x.not(), y), MisNegasiKeliru)
		for _, l := range []lit{x, x.not(), y, y.not()} {
			add(fact(l), MisKesimpulanTakTerkait)
		}
		res.steps = []string{
			fmt.Sprintf("Premis \"%s\" setara dengan kontraposisinya \"%s\"; premis \"%s\" setara dengan \"%s\".", p.propText(p.cprem[0]), p.propText(contra(p.cprem[0])), p.propText(p.cprem[1]), p.propText(contra(p.cprem[1]))),
			fmt.Sprintf("Rantaikan kedua premis (silogisme hipotetis): \"%s\" pasti benar, dan kontraposisinya \"%s\" juga.", p.propText(p.ckey), p.propText(contra(p.ckey))),
		}
	default:
		return result{}, false
	}
	return res, true
}

// chooseDistractors takes n wrong options, one per mistake in turn (mistake
// order and the options within it are shuffled), so the mistakes are mixed.
func chooseDistractors(rng *rand.Rand, cands []cand, n int) []cand {
	var order []string
	groups := map[string][]cand{}
	for _, c := range cands {
		if _, ok := groups[c.mis]; !ok {
			order = append(order, c.mis)
		}
		groups[c.mis] = append(groups[c.mis], c)
	}
	rng.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	for _, m := range order {
		g := groups[m]
		rng.Shuffle(len(g), func(i, j int) { g[i], g[j] = g[j], g[i] })
	}
	var out []cand
	for round := 0; len(out) < n; round++ {
		progressed := false
		for _, m := range order {
			if round < len(groups[m]) && len(out) < n {
				out = append(out, groups[m][round])
				progressed = true
			}
		}
		if !progressed {
			break
		}
	}
	return out
}

func build(p params, seed int64) gen.Draft {
	res, _ := solve(p)
	rng := rand.New(rand.NewSource(derive(seed, shuffleSalt)))
	type opt struct{ text, mis string }
	opts := []opt{{res.key, ""}}
	steps := append([]string(nil), res.steps...)
	warned := map[string]bool{}
	for _, c := range chooseDistractors(rng, res.cands, optionCount-1) {
		opts = append(opts, opt{c.text, c.mis})
		if !warned[c.mis] {
			warned[c.mis] = true
			steps = append(steps, fmt.Sprintf("Kesalahan umum: %s.", misconceptionLabels[c.mis]))
		}
	}
	rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	options := make([]gen.Option, len(opts))
	answerKey := ""
	for i, o := range opts {
		k := string(rune('A' + i))
		options[i] = gen.Option{Key: k, Text: o.text, MisconceptionID: o.mis}
		if o.mis == "" {
			answerKey = k
		}
	}
	return gen.Draft{
		PatternID: PatternID,
		Stem:      stem(p),
		Options:   options,
		AnswerKey: answerKey,
		Steps:     steps,
		Seed:      seed,
	}
}

func stem(p params) string {
	var p1, p2 string
	if p.kind == kindCat {
		p1, p2 = p.stmtText(p.prem[0]), p.stmtText(p.prem[1])
	} else {
		p1, p2 = p.propText(p.cprem[0]), p.propText(p.cprem[1])
	}
	return fmt.Sprintf("Premis 1: %s. Premis 2: %s. Kesimpulan yang tepat adalah ...", p1, p2)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
