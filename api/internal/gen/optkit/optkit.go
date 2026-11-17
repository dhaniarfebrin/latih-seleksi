// Package optkit holds the seed mixing and option assembly shared by the
// newer question families.
package optkit

import (
	"fmt"
	"math/rand"

	"latihseleksi/api/internal/gen"
)

const (
	// OptionCount is the number of options of every question.
	OptionCount = 5
	shuffleSalt = 1 << 20
)

// Derive mixes (seed, k) into a new seed (splitmix64).
func Derive(seed int64, k int) int64 {
	z := uint64(seed) + uint64(k+1)*0x9E3779B97F4A7C15
	z = (z ^ z>>30) * 0xBF58476D1CE4E5B9
	z = (z ^ z>>27) * 0x94D049BB133111EB
	return int64(z ^ z>>31)
}

// Cand is a wrong-option candidate with the misconception behind it.
type Cand struct {
	Text string
	Mis  string
}

// Build assembles a draft: the key plus the first four distinct candidates
// (in the order given), shuffled deterministically from seed. It errors when
// fewer than four distinct candidates differ from the key.
func Build(patternID, stem string, seed int64, key string, cands []Cand, steps []string) (gen.Draft, error) {
	opts := []Cand{{Text: key}}
	seen := map[string]bool{key: true}
	for _, c := range cands {
		if len(opts) == OptionCount {
			break
		}
		if c.Text == "" || seen[c.Text] {
			continue
		}
		seen[c.Text] = true
		opts = append(opts, c)
	}
	if len(opts) < OptionCount {
		return gen.Draft{}, fmt.Errorf("%s: hanya %d opsi berbeda", patternID, len(opts))
	}
	rng := rand.New(rand.NewSource(Derive(seed, shuffleSalt)))
	rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	options := make([]gen.Option, len(opts))
	answer := ""
	for i, o := range opts {
		k := string(rune('A' + i))
		options[i] = gen.Option{Key: k, Text: o.Text, MisconceptionID: o.Mis}
		if o.Mis == "" {
			answer = k
		}
	}
	return gen.Draft{PatternID: patternID, Stem: stem, Options: options, AnswerKey: answer, Steps: steps, Seed: seed}, nil
}

// Retry runs try with attempt-specific seeds until it returns ok, at most
// attempts times.
func Retry(name string, seed int64, band gen.Band, attempts int, try func(rng *rand.Rand) (gen.Draft, bool)) (gen.Draft, error) {
	if band < gen.Band1 || band > gen.Band3 {
		return gen.Draft{}, fmt.Errorf("%s: band %d tidak dikenal", name, band)
	}
	for k := 0; k < attempts; k++ {
		if d, ok := try(rand.New(rand.NewSource(Derive(seed, k)))); ok {
			d.Seed = seed
			return d, nil
		}
	}
	return gen.Draft{}, fmt.Errorf("%s: tidak ada parameter valid untuk seed %d band %d dalam %d percobaan", name, seed, band, attempts)
}
