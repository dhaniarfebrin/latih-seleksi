// Package gentest holds the checks every question family test shares.
package gentest

import (
	"reflect"
	"testing"

	"latihseleksi/api/internal/gen"
)

// Run generates seeds x 3 bands questions and checks the contract shared by
// all families (determinism, five distinct options, key among them, a
// misconception label per wrong option, steps), then hands each draft to
// check for the family's independent key test. It fails when a band yields
// fewer than minDistinct distinct stems.
func Run(t *testing.T, fam gen.Family, seeds, minDistinct int, check func(t *testing.T, d gen.Draft, band gen.Band)) {
	t.Helper()
	labels := fam.Misconceptions()
	for _, band := range []gen.Band{gen.Band1, gen.Band2, gen.Band3} {
		stems := map[string]bool{}
		for s := 0; s < seeds; s++ {
			seed := int64(s)*7919 + 13
			d, err := fam.Generate(seed, band)
			if err != nil {
				t.Fatalf("band %d seed %d: %v", band, seed, err)
			}
			again, _ := fam.Generate(seed, band)
			if !reflect.DeepEqual(d, again) {
				t.Fatalf("band %d seed %d: not deterministic", band, seed)
			}
			if d.PatternID != fam.ID() || d.Seed != seed || len(d.Steps) == 0 {
				t.Fatalf("band %d seed %d: bad pattern id, seed or steps", band, seed)
			}
			if len(d.Options) != 5 {
				t.Fatalf("band %d seed %d: %d options", band, seed, len(d.Options))
			}
			texts, keys, correct := map[string]bool{}, map[string]bool{}, 0
			for _, o := range d.Options {
				if o.Text == "" || texts[o.Text] || keys[o.Key] {
					t.Fatalf("band %d seed %d: duplicate or empty option in %q", band, seed, d.Stem)
				}
				texts[o.Text], keys[o.Key] = true, true
				if o.MisconceptionID == "" {
					correct++
					if o.Key != d.AnswerKey {
						t.Fatalf("band %d seed %d: correct option is not the answer key", band, seed)
					}
				} else if _, ok := labels[o.MisconceptionID]; !ok {
					t.Fatalf("band %d seed %d: misconception %q has no label", band, seed, o.MisconceptionID)
				}
			}
			if correct != 1 || !keys[d.AnswerKey] {
				t.Fatalf("band %d seed %d: want exactly one correct option", band, seed)
			}
			stems[d.Stem] = true
			check(t, d, band)
			if t.Failed() {
				t.Fatalf("band %d seed %d: %q", band, seed, d.Stem)
			}
		}
		if len(stems) < minDistinct {
			t.Errorf("band %d: %d distinct stems in %d seeds, want >= %d", band, len(stems), seeds, minDistinct)
		}
	}
}

// KeyText returns the text of the option the draft names as its answer.
func KeyText(d gen.Draft) string {
	for _, o := range d.Options {
		if o.Key == d.AnswerKey {
			return o.Text
		}
	}
	return ""
}
