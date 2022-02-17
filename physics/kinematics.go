package physics

import (
	"math/big"

	"git.sr.ht/~kisom/proxima/rat"
)

var (
	// C is the speed of light in m/s.
	C = rat.Float(299792458.0)

	// Gee is one standard earth gravitational acceleration in m/s^2.
	Gee = rat.Float(9.800665) // m/s^2

	// NegGee is a negative Gee.
	NegGee = rat.Float(-9.800665) // m/s^2
)

// State tracks an object's velocity and position.
type State struct {
	X *big.Rat
	V *big.Rat
}

// Accelerate computes the state after an acceleration period.
func (s State) Accelerate(a *big.Rat, seconds float64) State {
	t := rat.Float(seconds)
	return State{
		X: AccelerationDistance(s.X, s.V, a, t),
		V: AccelerationVelocity(s.V, a, t),
	}
}

// Coast computes the state after a period of coasting.
func (s State) Coast(seconds float64) State {
	t := rat.Float(seconds)
	return State{
		X: rat.Add(s.X, rat.Mul(s.V, t)),
		V: s.V,
	}
}

// AccelerationDistance computes the distance traveled during an acceleration.
func AccelerationDistance(x0, v0, a, t *big.Rat) *big.Rat {
	vt := rat.Mul(v0, t)
	at2 := rat.Mul(rat.Half, rat.Mul(a, rat.Sqr(t)))
	return rat.Add(x0, rat.Add(vt, at2))
}

// AccelerationVelocity computes the final velocity after an acceleration.
func AccelerationVelocity(v0, a, t *big.Rat) *big.Rat {
	return rat.Add(v0, rat.Mul(a, t))
}
