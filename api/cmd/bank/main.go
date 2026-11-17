// Command bank validates the embedded TWK/TKP bank and prints review lists.
package main

import (
	"fmt"
	"os"
	"sort"

	"latihseleksi/api/internal/bank"
)

const usage = "usage: bank check | bank review | bank sheet | bank sheet-twk"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	b, err := bank.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "check":
		check(b)
	case "review":
		review(b)
	case "sheet":
		sheet(b)
	case "sheet-twk":
		sheetTWK(b)
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
}

func check(b *bank.Bank) {
	fmt.Printf("OK: %d TWK, %d TKP, rubric %d values\n", len(b.TWK), len(b.TKP), len(b.Rubric.Values))
	for _, section := range []string{bank.SectionTWK, bank.SectionTKP} {
		cov := b.Coverage(section)
		aspects := make([]string, 0, len(cov))
		for a := range cov {
			aspects = append(aspects, a)
		}
		sort.Strings(aspects)
		fmt.Printf("%s coverage (reviewed/draft per aspek):\n", section)
		for _, a := range aspects {
			fmt.Printf("  %-18s %d/%d\n", a, cov[a].Reviewed, cov[a].Draft)
		}
	}
}

func review(b *bank.Bank) {
	for _, e := range b.TWK {
		if e.Status == bank.StatusDraft {
			fmt.Printf("[%s] %s (%s)\n  %s\n  kunci: %s | kutipan: %q\n", e.Aspek, e.ID, e.Ref, e.Question, e.Correct, e.Quote)
		}
	}
	for _, e := range b.TKP {
		if e.Status != bank.StatusDraft {
			continue
		}
		fmt.Printf("[%s] %s (rubrik %s/%s)\n  %s\n", e.Aspek, e.ID, e.RubricRef.Value, e.RubricRef.Butir, e.Scenario)
		for _, o := range e.Options {
			fmt.Printf("  bobot %d: %s\n", o.Weight, o.Text)
		}
	}
}

// sheet prints a markdown review sheet of every TKP item, for a human to
// approve the weights item by item.
func sheet(b *bank.Bank) {
	butir := map[bank.RubricRef]string{}
	label := map[string]string{}
	for _, v := range b.Rubric.Values {
		label[v.ID] = v.Label
		for _, x := range v.Butir {
			butir[bank.RubricRef{Value: v.ID, Butir: x.ID}] = x.Text
		}
	}
	fmt.Printf("# Lembar tinjau bobot TKP\n\nSumber rubrik: %s. Bobot adalah penilaian penyusun, bukan kunci resmi.\n", b.Rubric.Source)
	for _, e := range b.TKP {
		fmt.Printf("\n## %s\n\n- status: %s\n- aspek: %s\n- butir rubrik: %s/%s (%s): %s\n- skenario: %s\n- opsi:\n",
			e.ID, e.Status, e.Aspek, e.RubricRef.Value, e.RubricRef.Butir, label[e.RubricRef.Value], butir[e.RubricRef], e.Scenario)
		for _, o := range e.Options {
			fmt.Printf("  - bobot %d: %s\n", o.Weight, o.Text)
		}
	}
}

// sheetTWK prints a markdown review sheet of every TWK item, grouped by
// aspek, for a human to read the question, key and quote against the source.
func sheetTWK(b *bank.Bank) {
	fmt.Print("# Lembar tinjau TWK\n\nKunci berasal dari kutipan sumber (dicek otomatis terhadap korpus). Yang perlu dinilai manusia: kutipan benar dan utuh, kunci tidak ambigu, pengecoh jelas salah, kalimat soal wajar. Tandai tiap item: [x] setuju, atau tulis catatan.\n")
	byAspek := map[string][]bank.TWKEntry{}
	var aspects []string
	for _, e := range b.TWK {
		if _, ok := byAspek[e.Aspek]; !ok {
			aspects = append(aspects, e.Aspek)
		}
		byAspek[e.Aspek] = append(byAspek[e.Aspek], e)
	}
	sort.Strings(aspects)
	for _, a := range aspects {
		fmt.Printf("\n# Aspek: %s (%d item)\n", a, len(byAspek[a]))
		for _, e := range byAspek[a] {
			fmt.Printf("\n## %s\n\n- [ ] setuju | catatan:\n- status: %s, ditinjau: %s\n- sumber: %s, %s\n- soal: %s\n- kunci: %s\n- pengecoh:\n", e.ID, e.Status, e.ReviewedBy, e.SourceID, e.Ref, e.Question, e.Correct)
			for _, d := range e.Distractors {
				fmt.Printf("  - %s\n", d)
			}
			fmt.Printf("- kutipan: %q\n", e.Quote)
		}
	}
}
