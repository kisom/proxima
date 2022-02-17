package physics

import (
	"log"
	"math"
	"testing"
	"time"

	"git.sr.ht/~kisom/proxima/rat"
)

func TestLorentz(t *testing.T) {
	v := PercentCToVelocity(0.999)
	expected := rat.Float(0.0447102)

	γ := lorentz(v)
	Δ := rat.Sub(γ, expected)
	if σ, _ := Δ.Float64(); math.Abs(σ) > 0.00001 {
		t.Fatalf("Lorenz should return γ=%s, but have %s (σ=%0.4f)",
			expected.FloatString(4), γ.FloatString(4), σ)
	}
}

func TestRelativeTimeMinor(t *testing.T) {
	v := PercentCToVelocity(0.1)
	dt := time.Hour
	expected := 18 * time.Second // expected difference

	d := relativeTime(dt.Seconds(), v)
	delta := math.Abs(d.Seconds() - (dt.Seconds() + expected.Seconds()))
	if delta > 0.0015 {
		t.Fatalf("relative time should be ±0.0015, delta is %f", delta)
	}
}

func TestRelativeTimeNoticeable(t *testing.T) {
	// Testing how much time passes:
	v := PercentCToVelocity(0.999)
	Δt := time.Hour + 30*time.Minute

	// Traveling at v for Δt, the clock on the ship should pass by Δt,
	// while more time should pass by on Earth.
	//
	// Specifically, if 1.5h pass by on the ship at 99.9% of the speed of
	// light, then ~33h32m58s should pass by on Earth.
	expected := 33*time.Hour + 32*time.Minute + 58*time.Second

	earthTime := relativeTime(Δt.Seconds(), v)

	// Now, compare the difference between the actual and expected times.
	σ := math.Abs(expected.Seconds() - earthTime.Seconds())
	if σ > 0.001 {
		log.Fatalf("σ should be < 0.001s, but it is %0.4fs", σ)
	}
}

func TestClocks(t *testing.T) {
	v := PercentCToVelocity(0.999)
	Δt := time.Hour + 30*time.Minute

	// The difference between the clocks should be the relative time passed
	// less the change in time on the ship.
	expected := 33*time.Hour + 32*time.Minute + 58*time.Second - Δt

	clock := Now()
	clock.Update(Δt.Seconds(), v)

	δ := clock.Drift()
	σ := math.Abs(δ.Seconds() - expected.Seconds())
	if σ > 0.001 {
		log.Fatalf("σ should be < 0.001s, but it is %0.4fs", σ)
	}
}
