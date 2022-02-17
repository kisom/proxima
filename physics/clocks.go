package physics

import (
	"math/big"
	"time"

	"git.sr.ht/~kisom/proxima/rat"
	"github.com/benbjohnson/clock"
)

type Clock struct {
	Launch   time.Time
	Observer *clock.Mock
	Relative *clock.Mock
}

// Now returns a new Clock initialized to realtime now.
func Now() *Clock {
	now := time.Now().In(time.UTC)

	c := &Clock{
		Launch:   now,
		Observer: clock.NewMock(),
		Relative: clock.NewMock(),
	}
	c.Observer.Set(now)
	c.Relative.Set(now)
	return c
}

// Sync sets the observer clock to the system time.
func (c *Clock) Sync() {
	now := time.Now().In(time.UTC)
	c.Observer.Set(now)
}

func (c *Clock) Update(seconds float64, v *big.Rat) {
	dt := time.Duration(seconds * float64(time.Second))
	rt := relativeTime(seconds, v)
	c.Observer.Add(dt)
	c.Relative.Add(rt)
}

func (c *Clock) Drift() time.Duration {
	return c.Relative.Now().Sub(c.Observer.Now())
}

func lorentz(v *big.Rat) *big.Rat {
	gamma := rat.Div(rat.Sqr(v), rat.Sqr(C))
	gamma = rat.Sub(rat.One, gamma)
	return rat.Sqrt(gamma)
}

func relativeTime(dt float64, v *big.Rat) time.Duration {
	// Δt is in seconds.
	// Δt_relative = Δt / γ
	// γ = √(1 – v²/c²)
	gamma := lorentz(v)
	dtRat := rat.Float(dt)
	dtr := rat.Div(dtRat, gamma)
	dtr = rat.Mul(dtr, rat.Second)
	return rat.Duration(dtr).Round(time.Second)
}
