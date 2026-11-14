package kerjasama

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
	timeRe     = regexp.MustCompile(`dalam (\d+) hari`)
	togetherRe = regexp.MustCompile(`selama (\d+) hari`)
	optionRe   = regexp.MustCompile(`^(\d+) hari$`)
)

// independentDays recomputes the total days from the stem text only. Joint
// work uses the sum of rates; the "B stops" case uses the closed form
// a*(1 - k/b), which differs from the implementation's step-by-step path.
func independentDays(t *testing.T, stem string) *big.Rat {
	t.Helper()
	ms := timeRe.FindAllStringSubmatch(stem, -1)
	if len(ms) < 2 {
		t.Fatalf("stem has too few times: %q", stem)
	}
	times := make([]int64, len(ms))
	for i, m := range ms {
		v, _ := strconv.Atoi(m[1])
		times[i] = int64(v)
	}
	if tm := togetherRe.FindStringSubmatch(stem); tm != nil {
		if len(times) != 2 {
			t.Fatalf("band 3 stem should have two times: %q", stem)
		}
		k, _ := strconv.Atoi(tm[1])
		frac := new(big.Rat).Sub(big.NewRat(1, 1), big.NewRat(int64(k), times[1]))
		return frac.Mul(frac, big.NewRat(times[0], 1))
	}
	rate := new(big.Rat)
	for _, v := range times {
		rate.Add(rate, big.NewRat(1, v))
	}
	return rate.Inv(rate)
}

func optionDays(t *testing.T, text string) *big.Rat {
	t.Helper()
	m := optionRe.FindStringSubmatch(text)
	if m == nil {
		t.Fatalf("unparsable option %q", text)
	}
	v, _ := strconv.Atoi(m[1])
	return big.NewRat(int64(v), 1)
}

func TestPropertyKeyMatchesIndependentCalculation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			want := independentDays(t, d.Stem)
			if !want.IsInt() || want.Sign() <= 0 {
				t.Fatalf("band %d seed %d: key %s is not a positive whole number", band, seed, want.RatString())
			}

			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			matches := 0
			texts := map[string]bool{}
			for _, o := range d.Options {
				v := optionDays(t, o.Text)
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
			name: "band1 A 12 hari, B 24 hari",
			p:    params{times: []int{12, 24}},
			stem: "A dapat mengecat sebuah pagar sendirian dalam 12 hari, dan B dalam 24 hari. Jika semuanya bekerja bersama, pekerjaan itu selesai dalam ... hari.",
			key:  "8 hari",
			wantMis: map[string][]string{
				MisRataRataWaktu:    {"18 hari"},
				MisJumlahWaktu:      {"36 hari"},
				MisWaktuTercepat:    {"12 hari"},
				MisSalahHitungDekat: {"9 hari"},
			},
			stepPart: "Waktu bersama = 1 : 1/8 = 8 hari.",
		},
		{
			name: "band2 A 12, B 15, C 20 hari",
			p:    params{times: []int{12, 15, 20}},
			stem: "A dapat mengecat sebuah pagar sendirian dalam 12 hari, B dalam 15 hari, dan C dalam 20 hari. Jika semuanya bekerja bersama, pekerjaan itu selesai dalam ... hari.",
			key:  "5 hari",
			wantMis: map[string][]string{
				MisJumlahWaktu:      {"47 hari"},
				MisWaktuTercepat:    {"12 hari"},
				MisSalahHitungDekat: {"4 hari", "6 hari"},
			},
			stepPart: "Laju bersama = 1/12 + 1/15 + 1/20 = 1/5 pekerjaan per hari.",
		},
		{
			name: "band3 A 12, B 24, bersama 2 hari",
			p:    params{times: []int{12, 24}, together: 2},
			stem: "A dapat mengecat sebuah pagar sendirian dalam 12 hari, dan B dalam 24 hari. A dan B bekerja bersama selama 2 hari, kemudian B berhenti dan A melanjutkan sendirian sampai selesai. Seluruh pekerjaan itu selesai dalam ... hari.",
			key:  "11 hari",
			wantMis: map[string][]string{
				MisHanyaSisaHari:    {"9 hari"},
				MisWaktuSendiri:     {"12 hari"},
				MisSalahHitungDekat: {"10 hari", "13 hari"},
			},
			stepPart: "Total = 2 + 9 = 11 hari.",
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
	if fam.ID() != PatternID || fam.Label() != "Kerja sama" {
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
