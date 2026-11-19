package catalog

import "testing"

func TestSubset(t *testing.T) {
	got := Subset("TIU.verbal.silogisme", "TIU.numerik.deret-angka")
	if len(got) != 2 || got[0].ID() != "TIU.verbal.silogisme" || got[1].ID() != "TIU.numerik.deret-angka" {
		t.Fatalf("Subset order wrong: %v", got)
	}
	defer func() {
		if recover() == nil {
			t.Error("unknown id did not panic")
		}
	}()
	Subset("nope")
}
