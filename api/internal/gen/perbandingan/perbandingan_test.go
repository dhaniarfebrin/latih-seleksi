package perbandingan

import (
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
)

var (
	splitRe  = regexp.MustCompile(`^Perbandingan (\w+) (.+?) dan (.+?) adalah (\d+) : (\d+)\. Jika jumlah keduanya (\d+)[^,]*, (\w+) (.+?) adalah \.\.\.$`)
	threeRe  = regexp.MustCompile(`^Perbandingan (\w+) (.+?), (.+?), dan (.+?) adalah (\d+) : (\d+) : (\d+)\. Jika (.+), (.+) adalah \.\.\.$`)
	diffRe   = regexp.MustCompile(`^selisih \w+ (.+?) dan (.+?) adalah (\d+)`)
	partRe   = regexp.MustCompile(`^\w+ (.+?) adalah (\d+)`)
	askRe    = regexp.MustCompile(`^\w+ (.+)$`)
	chainRe  = regexp.MustCompile(`^Perbandingan (\w+) (.+?) dan (.+?) adalah (\d+) : (\d+), sedangkan perbandingan (\w+) (.+?) dan (.+?) adalah (\d+) : (\d+)\. Jika jumlah ketiganya (\d+)[^,]*, (\w+) (.+?) adalah \.\.\.$`)
	ageRe    = regexp.MustCompile(`^Perbandingan umur (.+?) dan (.+?) sekarang adalah (\d+) : (\d+)\. (\d+) tahun (lagi|yang lalu), perbandingan umur mereka (?:menjadi|adalah) (\d+) : (\d+)\. Umur (.+?) sekarang adalah \.\.\. tahun\.$`)
	countRe  = regexp.MustCompile(`^Perbandingan banyak (.+?) dan (.+?) adalah (\d+) : (\d+)\. Setelah masing-masing (ditambah|dikurangi) (\d+) [a-z]+, perbandingannya menjadi (\d+) : (\d+)\. Banyak (.+?) semula adalah \.\.\. [a-z]+\.$`)
	digitsRe = regexp.MustCompile(`\d+`)
	optionRe = regexp.MustCompile(`^(\d+)`)

	// only inverse-proportion stems contain one of these phrases
	inverseMarkers = []string{"selesai dalam", "cukup untuk", "waktu tempuh", "kosong dalam"}
)

func atoi(s string) int64 {
	v, _ := strconv.Atoi(s)
	return int64(v)
}

func rat(s string) *big.Rat { return big.NewRat(atoi(s), 1) }

// variant names the stem shape recognised from the text.
func variant(stem string) string {
	switch {
	case splitRe.MatchString(stem):
		return "split"
	case threeRe.MatchString(stem):
		return "three"
	case chainRe.MatchString(stem):
		return "chain"
	case ageRe.MatchString(stem), countRe.MatchString(stem):
		return "shift"
	case strings.Contains(stem, "Perbandingan"):
		return ""
	}
	for _, m := range inverseMarkers {
		if strings.Contains(stem, m) {
			return "inverse"
		}
	}
	return "direct"
}

// indexOf returns the position of name in items, failing the test when absent.
func indexOf(t *testing.T, items []string, name string) int {
	t.Helper()
	for i, it := range items {
		if it == name {
			return i
		}
	}
	t.Fatalf("%q is not one of %v", name, items)
	return -1
}

// independentKey recomputes the answer from the stem text only, with exact
// rationals and a different route than the generator's steps.
func independentKey(t *testing.T, stem string) *big.Rat {
	t.Helper()
	switch variant(stem) {
	case "split":
		m := splitRe.FindStringSubmatch(stem)
		if m[1] != m[7] {
			t.Fatalf("measure changes inside %q", stem)
		}
		w := []int64{atoi(m[4]), atoi(m[5])}
		part := w[indexOf(t, []string{m[2], m[3]}, m[8])]
		v := new(big.Rat).Mul(rat(m[6]), big.NewRat(part, 1))
		return v.Quo(v, big.NewRat(w[0]+w[1], 1))
	case "three":
		m := threeRe.FindStringSubmatch(stem)
		items := []string{m[2], m[3], m[4]}
		w := []int64{atoi(m[5]), atoi(m[6]), atoi(m[7])}
		sum := w[0] + w[1] + w[2]
		var unit *big.Rat
		if d := diffRe.FindStringSubmatch(m[8]); d != nil {
			diff := w[indexOf(t, items, d[1])] - w[indexOf(t, items, d[2])]
			if diff < 0 {
				diff = -diff
			}
			unit = new(big.Rat).Quo(rat(d[3]), big.NewRat(diff, 1))
		} else if g := partRe.FindStringSubmatch(m[8]); g != nil {
			unit = new(big.Rat).Quo(rat(g[2]), big.NewRat(w[indexOf(t, items, g[1])], 1))
		} else {
			t.Fatalf("unreadable given clause in %q", stem)
		}
		if m[9] == "jumlah ketiganya" {
			return unit.Mul(unit, big.NewRat(sum, 1))
		}
		return unit.Mul(unit, big.NewRat(w[indexOf(t, items, askRe.FindStringSubmatch(m[9])[1])], 1))
	case "chain":
		m := chainRe.FindStringSubmatch(stem)
		if m[2] == "" || m[3] != m[7] || m[1] != m[6] || m[1] != m[12] {
			t.Fatalf("chain items do not line up in %q", stem)
		}
		items := []string{m[2], m[3], m[8]}
		// A = x, B = x*q/p, C = B*s/r; the total fixes x
		b := big.NewRat(atoi(m[5]), atoi(m[4]))
		c := new(big.Rat).Mul(b, big.NewRat(atoi(m[10]), atoi(m[9])))
		parts := []*big.Rat{big.NewRat(1, 1), b, c}
		total := new(big.Rat).Add(parts[0], new(big.Rat).Add(parts[1], parts[2]))
		x := new(big.Rat).Quo(rat(m[11]), total)
		return x.Mul(x, parts[indexOf(t, items, m[13])])
	case "shift":
		var a, b, n, c, d int64
		var sign int64 = 1
		var items []string
		var ask string
		if m := ageRe.FindStringSubmatch(stem); m != nil {
			items, ask = []string{m[1], m[2]}, m[9]
			a, b, n, c, d = atoi(m[3]), atoi(m[4]), atoi(m[5]), atoi(m[7]), atoi(m[8])
			if m[6] == "yang lalu" {
				sign = -1
			}
		} else {
			m := countRe.FindStringSubmatch(stem)
			items, ask = []string{m[1], m[2]}, m[9]
			a, b, n, c, d = atoi(m[3]), atoi(m[4]), atoi(m[6]), atoi(m[7]), atoi(m[8])
			if m[5] == "dikurangi" {
				sign = -1
			}
		}
		// search the multiplier k with (ak + sn) : (bk + sn) = c : d
		var found []int64
		for k := int64(1); k <= 2000; k++ {
			x, y := a*k+sign*n, b*k+sign*n
			if x > 0 && y > 0 && x*d == y*c {
				found = append(found, k)
			}
		}
		if len(found) != 1 {
			t.Fatalf("expected one multiplier, got %v for %q", found, stem)
		}
		w := []int64{a, b}
		return big.NewRat(w[indexOf(t, items, ask)]*found[0], 1)
	case "inverse", "direct":
		nums := digitsRe.FindAllString(stem, -1)
		if len(nums) != 3 {
			t.Fatalf("want three numbers in %q, got %v", stem, nums)
		}
		x1, y1, x2 := rat(nums[0]), rat(nums[1]), rat(nums[2])
		if variant(stem) == "inverse" {
			v := new(big.Rat).Mul(x1, y1)
			return v.Quo(v, x2)
		}
		v := new(big.Rat).Mul(y1, x2)
		return v.Quo(v, x1)
	}
	t.Fatalf("unrecognised stem %q", stem)
	return nil
}

func optionValue(t *testing.T, text string) *big.Rat {
	t.Helper()
	m := optionRe.FindStringSubmatch(text)
	if m == nil {
		t.Fatalf("unparsable option %q", text)
	}
	return rat(m[1])
}

func TestPropertyKeyMatchesIndependentCalculation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			want := independentKey(t, d.Stem)
			if !want.IsInt() || want.Sign() <= 0 {
				t.Fatalf("band %d seed %d: key %s is not a positive whole number: %q", band, seed, want.RatString(), d.Stem)
			}
			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			matches := 0
			texts := map[string]bool{}
			for _, o := range d.Options {
				v := optionValue(t, o.Text)
				if texts[o.Text] {
					t.Fatalf("band %d seed %d: duplicate option %q", band, seed, o.Text)
				}
				texts[o.Text] = true
				if v.Cmp(want) == 0 {
					matches++
					if o.Key != d.AnswerKey {
						t.Fatalf("band %d seed %d: AnswerKey %s but value matches %s", band, seed, d.AnswerKey, o.Key)
					}
					if o.MisconceptionID != "" {
						t.Fatalf("band %d seed %d: key option has misconception", band, seed)
					}
				} else if o.MisconceptionID == "" {
					t.Fatalf("band %d seed %d: distractor %s has no misconception", band, seed, o.Key)
				}
			}
			if matches != 1 {
				t.Fatalf("band %d seed %d: %d options equal key %s", band, seed, matches, want.RatString())
			}
			if d.Seed != seed || d.PatternID != PatternID {
				t.Fatalf("band %d seed %d: bad metadata %d %s", band, seed, d.Seed, d.PatternID)
			}
		}
	}
}

func TestBandsEmitExpectedVariants(t *testing.T) {
	want := map[gen.Band][]string{
		gen.Band1: {"split"},
		gen.Band2: {"three", "inverse", "direct"},
		gen.Band3: {"shift", "chain"},
	}
	for band, variants := range want {
		seen := map[string]int{}
		for seed := int64(1); seed <= 500; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			seen[variant(d.Stem)]++
		}
		for _, v := range variants {
			if seen[v] == 0 {
				t.Errorf("band %d never emitted variant %q (seen %v)", band, v, seen)
			}
		}
		if len(seen) != len(variants) {
			t.Errorf("band %d emitted unexpected variants %v", band, seen)
		}
	}
}

func TestShiftEmitsBothDirectionsAndContexts(t *testing.T) {
	later, earlier, added, removed := 0, 0, 0, 0
	for seed := int64(1); seed <= 2000; seed++ {
		d, err := Family{}.Generate(seed, gen.Band3)
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case strings.Contains(d.Stem, "tahun lagi"):
			later++
		case strings.Contains(d.Stem, "yang lalu"):
			earlier++
		case strings.Contains(d.Stem, "ditambah"):
			added++
		case strings.Contains(d.Stem, "dikurangi"):
			removed++
		}
	}
	if later == 0 || earlier == 0 || added == 0 || removed == 0 {
		t.Errorf("later=%d earlier=%d added=%d removed=%d, want all > 0", later, earlier, added, removed)
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
		if len(stems) < 80 {
			t.Errorf("band %d: only %d distinct stems in 200 seeds", band, len(stems))
		}
	}
}

// byMisconception groups distractor texts by misconception id and returns the key text.
func byMisconception(d gen.Draft) (string, map[string][]string) {
	out := map[string][]string{}
	key := ""
	for _, o := range d.Options {
		if o.Key == d.AnswerKey {
			key = o.Text
			continue
		}
		out[o.MisconceptionID] = append(out[o.MisconceptionID], o.Text)
	}
	for _, v := range out {
		sort.Strings(v)
	}
	return key, out
}

func TestFixtures(t *testing.T) {
	cases := []struct {
		name     string
		p        params
		stem     string
		key      string
		wantMis  map[string][]string
		stepPart string
	}{
		{
			name: "split 3:5 dari 160",
			p:    params{kind: kindSplit, w: []int{3, 5}, given: givenTotal, gv: 160, ask: 1},
			stem: "Perbandingan banyak kelereng merah dan kelereng biru adalah 3 : 5. Jika jumlah keduanya 160 butir, banyak kelereng biru adalah ...",
			key:  "100 butir",
			wantMis: map[string][]string{
				MisBagiTanpaJumlahRasio: {"32 butir"},
				MisBagianTertukar:       {"60 butir"},
				MisSalahHitungDekat:     {"101 butir", "99 butir"},
			},
			stepPart: "Satu bagian = 160 : 8 = 20.",
		},
		{
			name: "tiga bagian dengan selisih",
			p:    params{kind: kindThree, w: []int{2, 3, 6}, given: givenDiff, gi: 0, gj: 2, gv: 20, ask: 1},
			stem: "Perbandingan banyak kelereng merah, kelereng biru, dan kelereng hijau adalah 2 : 3 : 6. Jika selisih banyak kelereng merah dan kelereng hijau adalah 20 butir, banyak kelereng biru adalah ...",
			key:  "15 butir",
			wantMis: map[string][]string{
				MisBagiTanpaJumlahRasio: {"30 butir"},
				MisBagianTertukar:       {"10 butir"},
				MisSalahHitungDekat:     {"14 butir", "16 butir"},
			},
			stepPart: "Satu bagian = 20 : 4 = 5.",
		},
		{
			name: "tiga bagian, satu bagian diketahui, tanya jumlah",
			p:    params{kind: kindThree, context: 1, w: []int{2, 3, 5}, given: givenPart, gi: 1, gv: 30, ask: askTotal},
			stem: "Perbandingan banyak siswa kelas A, siswa kelas B, dan siswa kelas C adalah 2 : 3 : 5. Jika banyak siswa kelas B adalah 30 orang, jumlah ketiganya adalah ...",
			key:  "100 orang",
			wantMis: map[string][]string{
				MisBagiTanpaJumlahRasio: {"30 orang"},
				MisSalahHitungDekat:     {"101 orang", "102 orang", "99 orang"},
			},
			stepPart: "Jumlah ketiganya = 10 x 10 = 100.",
		},
		{
			name: "berbalik nilai",
			p:    params{kind: kindInverse, x1: 6, y1: 12, x2: 9},
			stem: "6 mesin pengepak dapat menyelesaikan satu pesanan dalam 12 hari. Jika dipakai 9 mesin yang sama, pesanan itu selesai dalam ... hari.",
			key:  "8 hari",
			wantMis: map[string][]string{
				MisSenilaiBukanBerbalik: {"18 hari"},
				MisSalahHitungDekat:     {"10 hari", "7 hari", "9 hari"},
			},
			stepPart: "y2 = (6 x 12) : 9 = 72 : 9.",
		},
		{
			name: "senilai",
			p:    params{kind: kindInverse, context: 4, x1: 6, y1: 9, x2: 14},
			stem: "Untuk membuat 6 kue dibutuhkan 9 butir telur. Untuk membuat 14 kue dibutuhkan ... butir telur.",
			key:  "21 butir",
			wantMis: map[string][]string{
				MisSalahHitungDekat: {"19 butir", "20 butir", "22 butir", "23 butir"},
			},
			stepPart: "y2 = (9 x 14) : 6 = 126 : 6.",
		},
		{
			name: "umur beberapa tahun lagi",
			p:    params{kind: kindShift, context: 2, w: []int{4, 5}, cd: [2]int{6, 7}, n: 12, sign: 1, ask: 0},
			stem: "Perbandingan umur kakak dan adik sekarang adalah 4 : 5. 12 tahun lagi, perbandingan umur mereka menjadi 6 : 7. Umur kakak sekarang adalah ... tahun.",
			key:  "24 tahun",
			wantMis: map[string][]string{
				MisSelisihBukanRasio: {"12 tahun", "36 tahun"},
				MisBagianTertukar:    {"30 tahun"},
				MisSalahHitungDekat:  {"25 tahun"},
			},
			stepPart: "kedua ruas dikalikan -1 menjadi k x 2 = 12, sehingga k = 6.",
		},
		{
			name: "jumlah dikurangi sama",
			p:    params{kind: kindShift, context: 4, w: []int{7, 8}, cd: [2]int{10, 13}, n: 11, sign: -1, ask: 1},
			stem: "Perbandingan banyak buku di rak kiri dan buku di rak kanan adalah 7 : 8. Setelah masing-masing dikurangi 11 buku, perbandingannya menjadi 10 : 13. Banyak buku di rak kanan semula adalah ... buku.",
			key:  "24 buku",
			wantMis: map[string][]string{
				MisSelisihBukanRasio: {"13 buku", "35 buku"},
				MisBagianTertukar:    {"21 buku"},
				MisSalahHitungDekat:  {"25 buku"},
			},
			stepPart: "k x 11 = 33 dan k = 3.",
		},
		{
			name: "perbandingan berantai",
			p:    params{kind: kindChain, context: 1, w: []int{5, 6, 3, 5}, gv: 84, ask: 1},
			stem: "Perbandingan banyak siswa kelas A dan siswa kelas B adalah 5 : 6, sedangkan perbandingan banyak siswa kelas B dan siswa kelas C adalah 3 : 5. Jika jumlah ketiganya 84 orang, banyak siswa kelas B adalah ...",
			key:  "24 orang",
			wantMis: map[string][]string{
				MisBagianTertukar:   {"20 orang", "40 orang"},
				MisSalahHitungDekat: {"23 orang", "25 orang"},
			},
			stepPart: "= 15 : 18 : 30 = 5 : 6 : 10.",
		},
		{
			name: "perbandingan berantai, bagian tengah tidak disamakan",
			p:    params{kind: kindChain, context: 2, w: []int{1, 2, 3, 2}, gv: 65, ask: 2},
			stem: "Perbandingan berat beras dan jagung adalah 1 : 2, sedangkan perbandingan berat jagung dan kacang tanah adalah 3 : 2. Jika jumlah ketiganya 65 kg, berat kacang tanah adalah ...",
			key:  "20 kg",
			wantMis: map[string][]string{
				MisBagianTertukar:     {"15 kg", "30 kg"},
				MisTengahTakDisamakan: {"26 kg"},
				MisSalahHitungDekat:   {"21 kg"},
			},
			stepPart: "Jumlah perbandingan = 3 + 6 + 4 = 13, satu bagian = 65 : 13 = 5.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !valid(c.p) {
				t.Fatalf("fixture params are rejected by valid")
			}
			d := build(c.p, 7)
			if d.Stem != c.stem {
				t.Errorf("stem = %q", d.Stem)
			}
			key, got := byMisconception(d)
			if key != c.key {
				t.Errorf("key = %q, want %q", key, c.key)
			}
			if !reflect.DeepEqual(got, c.wantMis) {
				t.Errorf("distractors = %v, want %v", got, c.wantMis)
			}
			if !strings.Contains(strings.Join(d.Steps, "\n"), c.stepPart) {
				t.Errorf("steps %q lack %q", d.Steps, c.stepPart)
			}
			if want := independentKey(t, d.Stem); optionValue(t, key).Cmp(want) != 0 {
				t.Errorf("key %q differs from independent value %s", key, want.RatString())
			}
		})
	}
}

func TestSolveRejectsMalformedParams(t *testing.T) {
	cases := map[string]params{
		"split tidak habis dibagi":      {kind: kindSplit, w: []int{3, 5}, given: givenTotal, gv: 161, ask: 0},
		"split tanpa bobot":             {kind: kindSplit},
		"split menanyakan jumlah":       {kind: kindSplit, w: []int{3, 5}, given: givenTotal, gv: 160, ask: askTotal},
		"selisih bobot sama":            {kind: kindThree, w: []int{2, 2, 5}, given: givenDiff, gi: 0, gj: 1, gv: 10, ask: 2},
		"bagian diketahui ditanyakan":   {kind: kindThree, w: []int{2, 3, 5}, given: givenPart, gi: 1, gv: 30, ask: 1},
		"berbalik tidak habis dibagi":   {kind: kindInverse, x1: 6, y1: 12, x2: 5},
		"x sama":                        {kind: kindInverse, x1: 6, y1: 12, x2: 6},
		"umur sebelumnya negatif":       {kind: kindShift, context: 0, w: []int{4, 5}, cd: [2]int{6, 7}, n: 12, sign: -1},
		"rasio sesudah sama":            {kind: kindShift, context: 4, w: []int{4, 5}, cd: [2]int{4, 5}, n: 3, sign: 1},
		"rantai tidak habis dibagi":     {kind: kindChain, w: []int{5, 6, 3, 5}, gv: 85, ask: 0},
		"rantai dengan bobot nol":       {kind: kindChain, w: []int{5, 0, 3, 5}, gv: 84, ask: 0},
		"jenis tidak dikenal":           {kind: kind(99)},
		"umur ayah di bawah batas umur": {kind: kindShift, context: 0, w: []int{3, 1}, cd: [2]int{2, 1}, n: 3, sign: 1},
	}
	for name, p := range cases {
		if _, ok := solve(p); ok {
			t.Errorf("%s: solve accepted %+v", name, p)
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
	if fam.ID() != PatternID || fam.Label() != "Perbandingan" {
		t.Errorf("id/label = %q/%q", fam.ID(), fam.Label())
	}
}

func TestGenerateErrorsWithoutValidParams(t *testing.T) {
	if _, err := generate(1, gen.Band2, maxAttempts, func(params) bool { return false }); err == nil {
		t.Fatal("expected error when no params are accepted")
	}
	if _, err := (Family{}).Generate(1, gen.Band(9)); err == nil {
		t.Fatal("expected error for unknown band")
	}
}

func TestDistinctStemsPer20000(t *testing.T) {
	min := map[gen.Band]int{gen.Band1: 12500, gen.Band2: 16500, gen.Band3: 16500}
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
		if len(stems) < min[band] {
			t.Errorf("band %d: only %d distinct stems in 20000 seeds, want >= %d", band, len(stems), min[band])
		}
	}
}
