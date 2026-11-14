package kecepatan

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
	distanceRe = regexp.MustCompile(`(?:berjarak|berada) (\d+) km`)
	speedRe    = regexp.MustCompile(`kecepatan (\d+) km/jam`)
	delayRe    = regexp.MustCompile(`(\d+) jam setelah A berangkat,`)
	optionRe   = regexp.MustCompile(`^(\d+) jam$`)
)

func atoi(t *testing.T, s string) int64 {
	t.Helper()
	v, err := strconv.Atoi(s)
	if err != nil {
		t.Fatal(err)
	}
	return int64(v)
}

// independentHours recomputes the meeting time from the stem text only. The
// late-start case solves vA*T + vB*(T-h) = D for T, which differs from the
// implementation's "head start, then closing speed" path.
func independentHours(t *testing.T, stem string) *big.Rat {
	t.Helper()
	dm := distanceRe.FindStringSubmatch(stem)
	sm := speedRe.FindAllStringSubmatch(stem, -1)
	if dm == nil || len(sm) != 2 {
		t.Fatalf("cannot parse distance/speeds from %q", stem)
	}
	d, a, b := atoi(t, dm[1]), atoi(t, sm[0][1]), atoi(t, sm[1][1])
	switch {
	case strings.Contains(stem, "mengejar"):
		return big.NewRat(d, a-b)
	case strings.Contains(stem, "waktu yang sama"):
		return big.NewRat(d, a+b)
	}
	hm := delayRe.FindStringSubmatch(stem)
	if hm == nil {
		t.Fatalf("cannot parse delay from %q", stem)
	}
	return big.NewRat(d+b*atoi(t, hm[1]), a+b)
}

func optionHours(t *testing.T, text string) *big.Rat {
	t.Helper()
	m := optionRe.FindStringSubmatch(text)
	if m == nil {
		t.Fatalf("unparsable option %q", text)
	}
	return big.NewRat(atoi(t, m[1]), 1)
}

func TestPropertyKeyMatchesIndependentCalculation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			want := independentHours(t, d.Stem)
			if !want.IsInt() || want.Sign() <= 0 {
				t.Fatalf("band %d seed %d: key %s is not a positive whole number", band, seed, want.RatString())
			}

			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			matches := 0
			texts := map[string]bool{}
			for _, o := range d.Options {
				v := optionHours(t, o.Text)
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
			name: "band1 40 dan 60 km/jam, 300 km",
			p:    params{kind: opposite, distance: 300, vA: 40, vB: 60},
			stem: "Kota P dan kota Q berjarak 300 km. Mobil A berangkat dari P menuju Q dengan kecepatan 40 km/jam, sedangkan mobil B berangkat dari Q menuju P dengan kecepatan 60 km/jam. Keduanya berangkat pada waktu yang sama. Mereka bertemu setelah ... jam.",
			key:  "3 jam",
			wantMis: map[string][]string{
				MisJumlahSelisihTertukar: {"15 jam"},
				MisKecepatanTunggal:      {"5 jam"},
				MisRataRataKecepatan:     {"6 jam"},
				MisSalahHitungDekat:      {"4 jam"},
			},
			stepPart: "Waktu bertemu = 300 : 100 = 3 jam.",
		},
		{
			name: "band2 60 menyusul 40 km/jam, 100 km",
			p:    params{kind: catchUp, distance: 100, vA: 60, vB: 40},
			stem: "Mobil A melaju dengan kecepatan 60 km/jam mengejar mobil B yang berada 100 km di depannya dan melaju searah dengan kecepatan 40 km/jam. A menyusul B setelah ... jam.",
			key:  "5 jam",
			wantMis: map[string][]string{
				MisJumlahSelisihTertukar: {"1 jam"},
				MisRataRataKecepatan:     {"2 jam"},
				MisSalahHitungDekat:      {"4 jam", "6 jam"},
			},
			stepPart: "Waktu menyusul = 100 : 20 = 5 jam.",
		},
		{
			name: "band3 60 dan 40 km/jam, B berangkat 1 jam kemudian, 260 km",
			p:    params{kind: lateStart, distance: 260, vA: 60, vB: 40, delay: 1},
			stem: "Kota P dan kota Q berjarak 260 km. Mobil A berangkat dari P menuju Q dengan kecepatan 60 km/jam. 1 jam setelah A berangkat, mobil B berangkat dari Q menuju P dengan kecepatan 40 km/jam. Mereka bertemu ... jam setelah A berangkat.",
			key:  "3 jam",
			wantMis: map[string][]string{
				MisLupaSelang:            {"2 jam"},
				MisJumlahSelisihTertukar: {"13 jam"},
				MisSalahHitungDekat:      {"4 jam", "5 jam"},
			},
			stepPart: "Waktu sejak A berangkat = 1 + 2 = 3 jam.",
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
	if fam.ID() != PatternID || fam.Label() != "Kecepatan relatif" {
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
