// Command bank validates the embedded TWK/TKP bank and prints review lists.
package main

import (
	"fmt"
	"os"
	"sort"

	"latihseleksi/api/internal/bank"
)

const usage = "usage: bank check | bank review"

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
