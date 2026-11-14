package persentase

import (
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"latihseleksi/api/internal/gen"
)

var changeRe = regexp.MustCompile(`(dinaikkan|diturunkan) (\d+)%`)

// independentPercent recomputes the net change from the stem text only,
// using exact rationals, without calling any implementation helper.
func independentPercent(t *testing.T, stem string) *big.Rat {
	t.Helper()
	ms := changeRe.FindAllStringSubmatch(stem, -1)
	if len(ms) < 2 {
		t.Fatalf("stem has too few changes: %q", stem)
	}
	r := big.NewRat(1, 1)
	for _, m := range ms {
		p, _ := strconv.Atoi(m[2])
		if m[1] == "diturunkan" {
			p = -p
		}
		r.Mul(r, big.NewRat(int64(100+p), 100))
	}
	r.Sub(r, big.NewRat(1, 1))
	return r.Mul(r, big.NewRat(100, 1))
}

// optionValue parses "naik 1,4%", "turun 10%" or "tetap" into a rational.
func optionValue(t *testing.T, text string) *big.Rat {
	t.Helper()
	if text == "tetap" {
		return new(big.Rat)
	}
	f := strings.Fields(text)
	if len(f) != 2 || (f[0] != "naik" && f[0] != "turun") || !strings.HasSuffix(f[1], "%") {
		t.Fatalf("unparsable option %q", text)
	}
	v, ok := new(big.Rat).SetString(strings.Replace(strings.TrimSuffix(f[1], "%"), ",", ".", 1))
	if !ok {
		t.Fatalf("unparsable number in %q", text)
	}
	if f[0] == "turun" {
		v.Neg(v)
	}
	return v
}

func TestPropertyKeyMatchesIndependentCalculation(t *testing.T) {
	fam := Family{}
	for band := gen.Band1; band <= gen.Band3; band++ {
		for seed := int64(1); seed <= 5000; seed++ {
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			want := independentPercent(t, d.Stem)

			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			matches := 0
			texts := map[string]bool{}
			values := map[string]bool{}
			for _, o := range d.Options {
				v := optionValue(t, o.Text)
				if texts[o.Text] || values[v.String()] {
					t.Fatalf("band %d seed %d: duplicate option %q", band, seed, o.Text)
				}
				texts[o.Text], values[v.String()] = true, true

				isKey := v.Cmp(want) == 0
				if isKey {
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
				t.Fatalf("band %d seed %d: %d options equal key %s", band, seed, matches, want.FloatString(4))
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
		if len(stems) < 100 {
			t.Errorf("band %d: only %d distinct stems in 200 seeds", band, len(stems))
		}
	}
}

func TestFixtureTIU08(t *testing.T) {
	p := gen.Params{Changes: []int{20, -25}, Context: 0}
	d := build(p, 7)
	ans, steps := Family{}.Solve(p)
	if ans.Text != "turun 10%" {
		t.Fatalf("answer = %q", ans.Text)
	}
	if len(steps) == 0 {
		t.Fatal("no steps")
	}
	byMis := map[string]string{}
	for _, o := range d.Options {
		if o.Key == d.AnswerKey {
			if o.Text != "turun 10%" {
				t.Errorf("key text = %q", o.Text)
			}
			continue
		}
		byMis[o.MisconceptionID] = o.Text
	}
	if byMis[MisJumlahLangsung] != "turun 5%" {
		t.Errorf("jumlah-langsung = %q, want turun 5%%", byMis[MisJumlahLangsung])
	}
	if byMis[MisLangkahTerakhir] != "turun 25%" {
		t.Errorf("langkah-terakhir-saja = %q", byMis[MisLangkahTerakhir])
	}
	if byMis[MisTandaTerbalik] != "naik 10%" {
		t.Errorf("tanda-terbalik = %q", byMis[MisTandaTerbalik])
	}
	if !strings.Contains(d.Stem, "dinaikkan 20%") || !strings.Contains(d.Stem, "diturunkan 25% dari harga setelah kenaikan") {
		t.Errorf("stem = %q", d.Stem)
	}
}

func TestGenerateErrorsWithoutValidParams(t *testing.T) {
	_, err := generate(1, gen.Band2, maxAttempts, func(gen.Params) bool { return false })
	if err == nil {
		t.Fatal("expected error when no params are accepted")
	}
	if _, err := (Family{}).Generate(1, gen.Band(9)); err == nil {
		t.Fatal("expected error for unknown band")
	}
}
