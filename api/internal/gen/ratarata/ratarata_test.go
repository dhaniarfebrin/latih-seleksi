package ratarata

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
	groupRe  = regexp.MustCompile(`terdiri dari (\d+) [a-z]+ dengan rata-rata [a-z ]+? (\d+)`)
	baseRe   = regexp.MustCompile(`Rata-rata [a-z ]+ dari (\d+) [a-z]+ adalah (\d+)`)
	addRe    = regexp.MustCompile(`lagi dengan [a-z ]+? (\d+)[a-z ]* ikut bergabung`)
	removeRe = regexp.MustCompile(`dengan [a-z ]+? (\d+)[a-z ]* keluar`)
	targetRe = regexp.MustCompile(`agar rata-rata [a-z ]+ menjadi (\d+)`)
	optionRe = regexp.MustCompile(`^(\d+)`)
)

func atoi(s string) int64 {
	v, _ := strconv.Atoi(s)
	return int64(v)
}

// variant names the stem shape recognised from the text.
func variant(stem string) string {
	switch {
	case groupRe.MatchString(stem):
		return "groups"
	case addRe.MatchString(stem):
		return "add"
	case removeRe.MatchString(stem):
		return "remove"
	case targetRe.MatchString(stem):
		return "target"
	}
	return ""
}

// independentKey recomputes the answer from the stem text only, with exact
// rationals and the weighted-mean definition rather than the generator's steps.
func independentKey(t *testing.T, stem string) *big.Rat {
	t.Helper()
	switch variant(stem) {
	case "groups":
		ms := groupRe.FindAllStringSubmatch(stem, -1)
		if len(ms) < 2 {
			t.Fatalf("too few groups: %q", stem)
		}
		total, heads := new(big.Rat), new(big.Rat)
		for _, m := range ms {
			n, r := atoi(m[1]), atoi(m[2])
			total.Add(total, big.NewRat(n*r, 1))
			heads.Add(heads, big.NewRat(n, 1))
		}
		return total.Quo(total, heads)
	case "add", "remove", "target":
		b := baseRe.FindStringSubmatch(stem)
		if b == nil {
			t.Fatalf("no base average: %q", stem)
		}
		n, r := big.NewRat(atoi(b[1]), 1), big.NewRat(atoi(b[2]), 1)
		sum := new(big.Rat).Mul(n, r)
		one := big.NewRat(1, 1)
		if m := addRe.FindStringSubmatch(stem); m != nil {
			sum.Add(sum, big.NewRat(atoi(m[1]), 1))
			return sum.Quo(sum, new(big.Rat).Add(n, one))
		}
		if m := removeRe.FindStringSubmatch(stem); m != nil {
			sum.Sub(sum, big.NewRat(atoi(m[1]), 1))
			return sum.Quo(sum, new(big.Rat).Sub(n, one))
		}
		m := targetRe.FindStringSubmatch(stem)
		want := new(big.Rat).Mul(big.NewRat(atoi(m[1]), 1), new(big.Rat).Add(n, one))
		return want.Sub(want, sum)
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
	return big.NewRat(atoi(m[1]), 1)
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
				t.Fatalf("band %d seed %d: key %s is not a positive whole number", band, seed, want.RatString())
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
		gen.Band1: {"groups"},
		gen.Band2: {"add", "remove", "target"},
		gen.Band3: {"groups"},
	}
	for band, variants := range want {
		seen := map[string]int{}
		groupCount := map[int]bool{}
		for seed := int64(1); seed <= 500; seed++ {
			d, err := Family{}.Generate(seed, band)
			if err != nil {
				t.Fatal(err)
			}
			seen[variant(d.Stem)]++
			groupCount[len(groupRe.FindAllString(d.Stem, -1))] = true
		}
		for _, v := range variants {
			if seen[v] == 0 {
				t.Errorf("band %d never emitted variant %q (seen %v)", band, v, seen)
			}
		}
		if len(seen) != len(variants) {
			t.Errorf("band %d emitted unexpected variants %v", band, seen)
		}
		if band == gen.Band1 && !(groupCount[2] && len(groupCount) == 1) {
			t.Errorf("band 1 should always have two groups, got %v", groupCount)
		}
		if band == gen.Band3 && !(groupCount[3] && len(groupCount) == 1) {
			t.Errorf("band 3 should always have three groups, got %v", groupCount)
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
			name: "band1 20 siswa rata-rata 70 dan 30 siswa rata-rata 80",
			p:    params{kind: kindGroups, groups: []group{{20, 70}, {30, 80}}},
			stem: "Kelompok pertama terdiri dari 20 siswa dengan rata-rata nilai 70. Kelompok kedua terdiri dari 30 siswa dengan rata-rata nilai 80. Rata-rata nilai gabungan seluruh kelompok itu adalah ...",
			key:  "76",
			wantMis: map[string][]string{
				MisRataRataDuaRata:  {"75"},
				MisPembagiSalah:     {"1900"},
				MisSalahHitungDekat: {"77", "78"},
			},
			stepPart: "Rata-rata gabungan = 3800 : 50 = 76.",
		},
		{
			name: "band2 tambah",
			p:    params{kind: kindAdd, n: 5, r: 60, x: 30, context: 1},
			stem: "Rata-rata berat badan dari 5 orang adalah 60 kg. Jika satu orang lagi dengan berat badan 30 kg ikut bergabung, rata-rata berat badan yang baru adalah ...",
			key:  "55 kg",
			wantMis: map[string][]string{
				MisRataRataDuaRata:  {"45 kg"},
				MisPembagiSalah:     {"66 kg"},
				MisSalahHitungDekat: {"54 kg", "56 kg"},
			},
			stepPart: "Jumlah baru = 300 + 30 = 330, banyak data = 5 + 1 = 6.",
		},
		{
			name: "band2 keluar",
			p:    params{kind: kindRemove, n: 9, r: 70, x: 38, context: 2},
			stem: "Rata-rata gaji dari 9 karyawan adalah 70 juta rupiah. Jika satu karyawan dengan gaji 38 juta rupiah keluar, rata-rata gaji yang baru adalah ...",
			key:  "74 juta rupiah",
			wantMis: map[string][]string{
				MisSalahHitungDekat: {"72 juta rupiah", "73 juta rupiah", "75 juta rupiah", "76 juta rupiah"},
			},
			stepPart: "Jumlah baru = 630 - 38 = 592, banyak data = 9 - 1 = 8.",
		},
		{
			name: "band2 target",
			p:    params{kind: kindTarget, n: 9, r: 70, target: 73},
			stem: "Rata-rata nilai dari 9 siswa adalah 70. Satu siswa lagi akan ditambahkan agar rata-rata nilai menjadi 73. Besar nilai yang ditambahkan itu adalah ...",
			key:  "100",
			wantMis: map[string][]string{
				MisRataRataDuaRata:  {"76"},
				MisPembagiSalah:     {"27"},
				MisSelisihDibalik:   {"43"},
				MisSalahHitungDekat: {"101"},
			},
			stepPart: "Nilai yang ditambahkan = 730 - 630 = 100.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
		})
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
	if fam.ID() != PatternID || fam.Label() != "Rata-rata gabungan" {
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
