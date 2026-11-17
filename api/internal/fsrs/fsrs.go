// Package fsrs schedules reviews with FSRS-6 using the default parameters.
//
// Formulas and the 21 default weights follow py-fsrs fsrs/scheduler.py (v6.3.3,
// raw file https://raw.githubusercontent.com/open-spaced-repetition/py-fsrs/main/fsrs/scheduler.py,
// read 2026-11-17). Differences from py-fsrs, on purpose: no learning or
// relearning steps (the next due date is always a whole number of days, at
// least 1) and no interval fuzz, so a review is a pure function of its input.
// Only Again and Good are used by the app; Hard and Easy are not modelled.
package fsrs

import (
	"math"
	"time"
)

// Grade is the outcome of one review.
type Grade int

// Values match py-fsrs Rating (Hard = 2 and Easy = 4 are unused here).
const (
	Again Grade = 1
	Good  Grade = 3
)

// DesiredRetention is the retrievability a card is rescheduled for.
const DesiredRetention = 0.9

const (
	minStability  = 0.001
	minDifficulty = 1.0
	maxDifficulty = 10.0
	maxInterval   = 36500
	hoursPerDay   = 24
)

// Weights are the FSRS-6 default parameters; the last one is the decay.
var Weights = [21]float64{
	0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001, 1.8722, 0.1666,
	0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014, 1.8729, 0.5425, 0.0912, 0.0658,
	0.1542,
}

var (
	decay  = -Weights[20]
	factor = math.Pow(DesiredRetention, 1/decay) - 1
)

// Card is the memory state of one item. A zero Reps means never reviewed.
type Card struct {
	Stability  float64
	Difficulty float64
	Due        time.Time
	LastReview time.Time
	Reps       int
	Lapses     int
}

// Retrievability is the predicted chance of recall at now; 0 for a new card.
func Retrievability(c Card, now time.Time) float64 {
	if c.Reps == 0 {
		return 0
	}
	return math.Pow(1+factor*float64(elapsedDays(c, now))/c.Stability, decay)
}

// Interval is the number of days after which retrievability falls to
// DesiredRetention, rounded, at least 1.
func Interval(stability float64) int {
	n := math.Round(stability / factor * (math.Pow(DesiredRetention, 1/decay) - 1))
	return int(math.Min(math.Max(n, 1), maxInterval))
}

// Review returns the card after a review with grade g at now.
func Review(c Card, g Grade, now time.Time) Card {
	if c.Reps == 0 {
		c.Stability = math.Max(Weights[g-1], minStability)
		c.Difficulty = clampDifficulty(initialDifficulty(g))
	} else {
		if elapsedDays(c, now) < 1 {
			c.Stability = shortTermStability(c.Stability, g)
		} else {
			c.Stability = nextStability(c.Difficulty, c.Stability, Retrievability(c, now), g)
		}
		c.Difficulty = nextDifficulty(c.Difficulty, g)
	}
	c.Reps++
	if g == Again {
		c.Lapses++
	}
	c.LastReview = now
	c.Due = now.Add(time.Duration(Interval(c.Stability)) * hoursPerDay * time.Hour)
	return c
}

// elapsedDays counts whole days since the last review, never negative.
func elapsedDays(c Card, now time.Time) int {
	d := int(now.Sub(c.LastReview).Hours() / hoursPerDay)
	if d < 0 {
		return 0
	}
	return d
}

func clampDifficulty(d float64) float64 { return math.Min(math.Max(d, minDifficulty), maxDifficulty) }

func initialDifficulty(g Grade) float64 {
	return Weights[4] - math.Exp(Weights[5]*float64(g-1)) + 1
}

func nextDifficulty(d float64, g Grade) float64 {
	delta := -(Weights[6] * float64(g-3))
	damped := d + (10-d)*delta/9
	return clampDifficulty(Weights[7]*initialDifficulty(4) + (1-Weights[7])*damped)
}

func shortTermStability(s float64, g Grade) float64 {
	inc := math.Exp(Weights[17]*(float64(g)-3+Weights[18])) * math.Pow(s, -Weights[19])
	if g != Again {
		inc = math.Max(inc, 1)
	}
	return math.Max(s*inc, minStability)
}

func nextStability(d, s, r float64, g Grade) float64 {
	var next float64
	if g == Again {
		long := Weights[11] * math.Pow(d, -Weights[12]) * (math.Pow(s+1, Weights[13]) - 1) * math.Exp((1-r)*Weights[14])
		short := s / math.Exp(Weights[17]*Weights[18])
		next = math.Min(long, short)
	} else {
		next = s * (1 + math.Exp(Weights[8])*(11-d)*math.Pow(s, -Weights[9])*(math.Exp((1-r)*Weights[10])-1))
	}
	return math.Max(next, minStability)
}
