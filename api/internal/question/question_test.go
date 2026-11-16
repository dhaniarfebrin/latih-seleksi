package question

import (
	"testing"

	"latihseleksi/api/internal/gen"
)

func TestValidAnswer(t *testing.T) {
	cases := map[string]bool{"A": true, "E": true, "C": true, "F": false, "a": false, "": false, "AB": false, "@": false}
	for in, want := range cases {
		if got := ValidAnswer(in); got != want {
			t.Errorf("ValidAnswer(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestRandomSeedNonNegative(t *testing.T) {
	for i := 0; i < 100; i++ {
		s, err := RandomSeed()
		if err != nil || s < 0 {
			t.Fatalf("RandomSeed() = %d, %v", s, err)
		}
	}
}

func TestFromDraftAndOption(t *testing.T) {
	d := gen.Draft{
		Stem: "q", AnswerKey: "B", Steps: []string{"s"},
		Options: []gen.Option{{Key: "A", Text: "x", MisconceptionID: "m1", Weight: 3}, {Key: "B", Text: "y"}},
	}
	sn := FromDraft(d)
	if sn.Stem != "q" || sn.AnswerKey != "B" || len(sn.Options) != 2 || sn.Options[0].MisconceptionID != "m1" || sn.Options[0].Weight != 3 || sn.Options[1].Weight != 0 {
		t.Fatalf("snapshot = %+v", sn)
	}
	if o, ok := sn.Option("A"); !ok || o.Text != "x" {
		t.Errorf("Option(A) = %+v, %v", o, ok)
	}
	if _, ok := sn.Option("Z"); ok {
		t.Error("Option(Z) found")
	}
}
