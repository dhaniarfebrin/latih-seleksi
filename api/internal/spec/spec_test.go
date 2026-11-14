package spec

import (
	"strings"
	"testing"
	"testing/fstest"
)

func mustLoad(t *testing.T) *Registry {
	t.Helper()
	r, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return r
}

func TestEmbeddedSpecsValidate(t *testing.T) {
	r := mustLoad(t)
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("got %d specs, want 2", len(all))
	}
	if all[0].ID != "skd-2024" || all[1].ID != "skd-sekolah-kedinasan-2026" {
		t.Errorf("All order = %s, %s", all[0].ID, all[1].ID)
	}
	for _, s := range all {
		if err := s.Validate(); err != nil {
			t.Errorf("%s: %v", s.ID, err)
		}
	}
	if _, ok := r.Get("nope"); ok {
		t.Error("Get(nope) found a spec")
	}
}

// sectionFacts checks the parts shared by both decrees: 110 items split
// 30/35/45 (321/2024 Diktum KEENAM, 406/2026 Diktum KEEMPAT), 100 minutes
// (Diktum KETIGA in both), maximum 550 = 150+175+225 (321/2024 Diktum
// KEDELAPAN, 406/2026 Diktum KEENAM), TWK 5 and TKP 6 aspects (Diktum KEDUA).
func sectionFacts(t *testing.T, s Spec) {
	t.Helper()
	if s.DurationMin != 100 {
		t.Errorf("%s duration = %d, want 100", s.ID, s.DurationMin)
	}
	want := []struct {
		id      string
		items   int
		max     int
		aspects int
	}{{"TWK", 30, 150, 5}, {"TIU", 35, 175, 0}, {"TKP", 45, 225, 6}}
	if len(s.Sections) != len(want) {
		t.Fatalf("%s: %d sections", s.ID, len(s.Sections))
	}
	items := 0
	for i, w := range want {
		sec := s.Sections[i]
		if sec.ID != w.id || sec.Items != w.items || sec.MaxScore() != w.max || len(sec.Aspects) != w.aspects {
			t.Errorf("%s %s: items=%d max=%d aspects=%d, want %d/%d/%d", s.ID, sec.ID, sec.Items, sec.MaxScore(), len(sec.Aspects), w.items, w.max, w.aspects)
		}
		items += sec.Items
	}
	if items != 110 || s.MaxTotal != 550 {
		t.Errorf("%s: items=%d max_total=%d, want 110/550", s.ID, items, s.MaxTotal)
	}
	// Diktum KETUJUH of 321/2024: TWK/TIU correct 5, wrong and blank 0; TKP 1..5, blank 0.
	for _, id := range []string{"TWK", "TIU"} {
		sc, _ := s.Section(id)
		if sc.Scoring != (Scoring{Kind: KindRightWrong, Correct: 5}) {
			t.Errorf("%s %s scoring = %+v", s.ID, id, sc.Scoring)
		}
	}
	tkp, _ := s.Section("TKP")
	if tkp.Scoring != (Scoring{Kind: KindWeighted, Min: 1, Max: 5}) {
		t.Errorf("%s TKP scoring = %+v", s.ID, tkp.Scoring)
	}
}

func trackFacts(t *testing.T, s Spec, id string, total int, th map[string]int) {
	t.Helper()
	tr, ok := s.Track(id)
	if !ok {
		t.Fatalf("%s: no track %s", s.ID, id)
	}
	if tr.TotalMin != total || len(tr.Thresholds) != len(th) {
		t.Errorf("%s/%s: total_min=%d thresholds=%v, want %d %v", s.ID, id, tr.TotalMin, tr.Thresholds, total, th)
	}
	for k, v := range th {
		if tr.Thresholds[k] != v {
			t.Errorf("%s/%s: %s = %d, want %d", s.ID, id, k, tr.Thresholds[k], v)
		}
	}
}

func TestGolden321_2024(t *testing.T) {
	s, ok := mustLoad(t).Get("skd-2024")
	if !ok {
		t.Fatal("skd-2024 missing")
	}
	if s.Scope != "cpns" || s.Source != "KepmenPANRB 321/2024" {
		t.Errorf("scope/source = %s / %s", s.Scope, s.Source)
	}
	sectionFacts(t, s)
	if len(s.Tracks) != 3 {
		t.Errorf("%d tracks, want 3", len(s.Tracks))
	}
	// Diktum KEDUA BELAS-KETIGA BELAS: umum 65/80/166.
	trackFacts(t, s, "umum", 0, map[string]int{"TWK": 65, "TIU": 80, "TKP": 166})
	// Diktum KELIMA BELAS: cumlaude dan diaspora, kumulatif 311 dan TIU 85.
	trackFacts(t, s, "cumlaude", 311, map[string]int{"TIU": 85})
	// Diktum KEENAM BELAS: disabilitas, Papua, daerah tertinggal, kumulatif 286 dan TIU 60.
	trackFacts(t, s, "afirmasi", 286, map[string]int{"TIU": 60})
}

func TestGolden406_2026(t *testing.T) {
	s, ok := mustLoad(t).Get("skd-sekolah-kedinasan-2026")
	if !ok {
		t.Fatal("skd-sekolah-kedinasan-2026 missing")
	}
	if s.Scope != "sekolah-kedinasan" || s.Source != "KepmenPANRB 406/2026" {
		t.Errorf("scope/source = %s / %s", s.Scope, s.Source)
	}
	sectionFacts(t, s)
	if len(s.Tracks) != 2 {
		t.Errorf("%d tracks, want 2", len(s.Tracks))
	}
	// Diktum KEDELAPAN: umum 65/80/156 (TKP differs from 166 of 321/2024).
	trackFacts(t, s, "umum", 0, map[string]int{"TWK": 65, "TIU": 80, "TKP": 156})
	// Diktum KESEPULUH: afirmasi daerah tertentu, kumulatif 281 dan TIU 55.
	trackFacts(t, s, "afirmasi", 281, map[string]int{"TIU": 55})
}

func TestValidateRejects(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(s *Spec)
		want   string
	}{
		{"threshold above section max", func(s *Spec) { s.Tracks[0].Thresholds["TWK"] = 151 }, "threshold TWK=151"},
		{"threshold zero", func(s *Spec) { s.Tracks[0].Thresholds["TWK"] = 0 }, "threshold TWK=0"},
		{"wrong max_total", func(s *Spec) { s.MaxTotal = 549 }, "max_total is 549"},
		{"duplicate track id", func(s *Spec) { s.Tracks[1].ID = s.Tracks[0].ID }, "duplicate track id"},
		{"duplicate section id", func(s *Spec) { s.Sections[1].ID = "TWK" }, "duplicate section id"},
		{"unknown section in threshold", func(s *Spec) { s.Tracks[0].Thresholds["XYZ"] = 10 }, "unknown section"},
		{"weighted min above max", func(s *Spec) { s.Sections[2].Scoring.Min = 6 }, "1 <= min <= max"},
		{"weighted min zero", func(s *Spec) { s.Sections[2].Scoring.Min = 0 }, "1 <= min <= max"},
		{"right-wrong correct zero", func(s *Spec) { s.Sections[0].Scoring.Correct = 0 }, "correct must be > 0"},
		{"negative wrong", func(s *Spec) { s.Sections[0].Scoring.Wrong = -1 }, "wrong and blank"},
		{"unknown scoring kind", func(s *Spec) { s.Sections[0].Scoring.Kind = "x" }, "unknown scoring kind"},
		{"max on right-wrong", func(s *Spec) { s.Sections[0].Scoring.Max = 5 }, "not used"},
		{"zero items", func(s *Spec) { s.Sections[0].Items = 0 }, "items must be > 0"},
		{"zero duration", func(s *Spec) { s.DurationMin = 0 }, "duration_min"},
		{"empty id", func(s *Spec) { s.ID = "" }, "id is empty"},
		{"empty source", func(s *Spec) { s.Source = "" }, "source is empty"},
		{"bad date", func(s *Spec) { s.EnactedOn = "29-07-2024" }, "enacted_on"},
		{"empty aspect", func(s *Spec) { s.Sections[0].Aspects[0] = "" }, "empty aspect"},
		{"duplicate aspect", func(s *Spec) { s.Sections[0].Aspects[1] = s.Sections[0].Aspects[0] }, "duplicate aspect"},
		{"total_min above max_total", func(s *Spec) { s.Tracks[1].TotalMin = 551 }, "total_min 551"},
		{"track without rule", func(s *Spec) { s.Tracks[1].TotalMin = 0; s.Tracks[1].Thresholds = nil }, "no rule"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, _ := mustLoad(t).Get("skd-2024") // Get returns a copy of the struct, but maps and slices are shared
			s = cloneSpec(s)
			c.mutate(&s)
			err := s.Validate()
			if err == nil {
				t.Fatal("Validate accepted an invalid spec")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q does not mention %q", err, c.want)
			}
		})
	}
}

// cloneSpec deep-copies the slices and maps so mutations do not leak between cases.
func cloneSpec(s Spec) Spec {
	out := s
	out.Sections = make([]Section, len(s.Sections))
	for i, sec := range s.Sections {
		sec.Aspects = append([]string(nil), sec.Aspects...)
		out.Sections[i] = sec
	}
	out.Tracks = make([]Track, len(s.Tracks))
	for i, tr := range s.Tracks {
		th := make(map[string]int, len(tr.Thresholds))
		for k, v := range tr.Thresholds {
			th[k] = v
		}
		tr.Thresholds = th
		out.Tracks[i] = tr
	}
	return out
}

func TestParseRejects(t *testing.T) {
	const base = `{"id":"x","label":"x","source":"x","source_url":"","enacted_on":"2024-01-01","scope":"x","duration_min":1,"sections":[],"tracks":[],"max_total":0`
	if _, err := Parse([]byte(base + `}`)); err != nil {
		t.Fatalf("base JSON should parse: %v", err)
	}
	if _, err := Parse([]byte(base + `,"surprise":1}`)); err == nil {
		t.Error("unknown field accepted")
	}
	if _, err := Parse([]byte(base + `} {}`)); err == nil {
		t.Error("trailing data accepted")
	}
	if _, err := Parse([]byte(`{`)); err == nil {
		t.Error("truncated JSON accepted")
	}
}

func TestLoadFSRejects(t *testing.T) {
	valid := `{"id":"a","label":"a","source":"a","source_url":"","enacted_on":"2024-01-01","scope":"x","duration_min":1,` +
		`"sections":[{"id":"S","label":"S","items":2,"scoring":{"kind":"right-wrong","correct":5}}],` +
		`"tracks":[{"id":"t","label":"t","thresholds":{"S":5}}],"max_total":10}`
	ok := fstest.MapFS{"d/a.json": {Data: []byte(valid)}}
	if _, err := loadFS(ok, "d"); err != nil {
		t.Fatalf("valid spec rejected: %v", err)
	}
	dup := fstest.MapFS{"d/a.json": {Data: []byte(valid)}, "d/b.json": {Data: []byte(valid)}}
	if _, err := loadFS(dup, "d"); err == nil || !strings.Contains(err.Error(), "duplicate spec id") {
		t.Errorf("duplicate id: err = %v", err)
	}
	bad := fstest.MapFS{"d/a.json": {Data: []byte(strings.Replace(valid, `"max_total":10`, `"max_total":11`, 1))}}
	if _, err := loadFS(bad, "d"); err == nil || !strings.Contains(err.Error(), "max_total") {
		t.Errorf("invalid spec: err = %v", err)
	}
	unk := fstest.MapFS{"d/a.json": {Data: []byte(strings.Replace(valid, `"scope":"x"`, `"scope":"x","extra":1`, 1))}}
	if _, err := loadFS(unk, "d"); err == nil {
		t.Error("unknown field accepted by loadFS")
	}
}
