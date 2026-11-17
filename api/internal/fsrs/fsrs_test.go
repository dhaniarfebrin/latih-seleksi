package fsrs

import (
	"math"
	"testing"
	"time"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func near(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %.12f, want %.12f", name, got, want)
	}
}

func TestWeightsCount(t *testing.T) {
	if len(Weights) != 21 {
		t.Fatalf("weights = %d, want 21 (FSRS-6)", len(Weights))
	}
}

// First Good review recomputed from the formulas by hand: stability is the
// third weight, difficulty = w4 - e^(w5*2) + 1.
func TestFirstReviewByHand(t *testing.T) {
	c := Review(Card{}, Good, t0)
	near(t, "stability", c.Stability, 2.3065)
	near(t, "difficulty", c.Difficulty, 6.4133-math.Exp(0.8334*2)+1)
	if c.Reps != 1 || c.Lapses != 0 || !c.LastReview.Equal(t0) {
		t.Errorf("card = %+v", c)
	}
	if got := c.Due.Sub(t0); got != 2*24*time.Hour {
		t.Errorf("due after %v, want 48h", got)
	}
	a := Review(Card{}, Again, t0)
	near(t, "again stability", a.Stability, 0.212)
	if a.Lapses != 1 {
		t.Errorf("lapses = %d", a.Lapses)
	}
}

// Second review by hand: R after 3 days, then the recall formula.
func TestSecondReviewByHand(t *testing.T) {
	c := Review(Card{}, Good, t0)
	now := t0.Add(3 * 24 * time.Hour)
	decay := -0.1542
	f := math.Pow(0.9, 1/decay) - 1
	r := math.Pow(1+f*3/c.Stability, decay)
	near(t, "retrievability", Retrievability(c, now), r)
	s, d := c.Stability, c.Difficulty
	want := s * (1 + math.Exp(1.8722)*(11-d)*math.Pow(s, -0.1666)*(math.Exp((1-r)*0.796)-1))
	got := Review(c, Good, now)
	near(t, "stability", got.Stability, want)
}

// Cross-check against py-fsrs 6.3.3 (Scheduler with learning_steps=(),
// relearning_steps=(), enable_fuzzing=False) for a fixed sequence.
func TestMatchesPyFSRS(t *testing.T) {
	steps := []struct {
		days     int
		grade    Grade
		stab     float64
		diff     float64
		interval int
	}{
		{0, Good, 2.3065, 2.118103970459016, 2},
		{3, Good, 13.826903694354568, 2.111214235785395, 14},
		{10, Again, 1.6661407475790257, 7.392238132342694, 2},
		{2, Good, 5.045934306639896, 7.38007426350719, 5},
		{0, Good, 5.045934306639896, 7.3679225585405215, 5},
		{30, Good, 25.649973707190163, 7.35578300527882, 26},
		{0, Again, 7.356083419352133, 9.116092223807549, 7},
	}
	var c Card
	now := t0
	for i, s := range steps {
		now = now.Add(time.Duration(s.days) * 24 * time.Hour)
		c = Review(c, s.grade, now)
		near(t, "stability", c.Stability, s.stab)
		near(t, "difficulty", c.Difficulty, s.diff)
		if got := int(c.Due.Sub(now).Hours() / 24); got != s.interval {
			t.Errorf("step %d: interval = %d, want %d", i, got, s.interval)
		}
	}
}

func TestAgainLowersStability(t *testing.T) {
	c := Review(Card{}, Good, t0)
	c = Review(c, Good, t0.Add(3*24*time.Hour))
	after := Review(c, Again, t0.Add(10*24*time.Hour))
	if after.Stability >= c.Stability || after.Lapses != 1 {
		t.Errorf("stability %v -> %v, lapses %d", c.Stability, after.Stability, after.Lapses)
	}
}

func TestIntervalAtDesiredRetentionEqualsStability(t *testing.T) {
	for _, s := range []float64{1, 2.3065, 13.8, 100, 365} {
		if got := Interval(s); got != int(math.Round(s)) {
			t.Errorf("Interval(%v) = %d", s, got)
		}
	}
	if Interval(0.001) != 1 || Interval(1e9) != maxInterval {
		t.Error("interval not clamped to [1, 36500]")
	}
}

func TestRetrievabilityOfNewCardIsZero(t *testing.T) {
	if Retrievability(Card{}, t0) != 0 {
		t.Error("new card has retrievability")
	}
}
