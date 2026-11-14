// Command content prints generated practice questions.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"latihseleksi/api/internal/gen"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
)

func main() {
	n := flag.Int("n", 8, "number of questions")
	seed := flag.Int64("seed", 1, "base seed; question i uses seed+i")
	band := flag.Int("band", 2, "difficulty band (1-3)")
	key := flag.Bool("key", false, "print answer key, steps and misconception labels")
	pattern := flag.String("pattern", persentase.PatternID, "pattern id")
	flag.Parse()

	var fam gen.Family
	var ids []string
	for _, f := range catalog.All() {
		ids = append(ids, f.ID())
		if f.ID() == *pattern {
			fam = f
		}
	}
	if fam == nil {
		fmt.Fprintf(os.Stderr, "error: unknown pattern %q; valid: %s\n", *pattern, strings.Join(ids, ", "))
		os.Exit(1)
	}
	for i := 0; i < *n; i++ {
		d, err := fam.Generate(*seed+int64(i), gen.Band(*band))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("%d. %s\n", i+1, d.Stem)
		for _, o := range d.Options {
			fmt.Printf("   %s. %s\n", o.Key, o.Text)
		}
		if *key {
			fmt.Printf("   Kunci: %s (seed %d)\n", d.AnswerKey, d.Seed)
			for _, s := range d.Steps {
				fmt.Printf("   - %s\n", s)
			}
			for _, o := range d.Options {
				if o.MisconceptionID != "" {
					fmt.Printf("   %s -> %s\n", o.Key, o.MisconceptionID)
				}
			}
		}
		fmt.Println()
	}
}
