package physics

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"git.sr.ht/~kisom/proxima/rat"
)

var (
	lyInMeters  = rat.Float(9.46073047258e+15)
	oneAU       = rat.Int64(149597870691)
	hundredthLY = LightyearsToMeters(0.01)
	hundredthAU = AstronomicalUnit(0.01)
	tenthAU     = AstronomicalUnit(0.1)
)

// LightyearsToMeters convert lightyears to meters.
func LightyearsToMeters(ly float64) *big.Rat {
	lyf := rat.Float(ly)
	return rat.Mul(lyf, lyInMeters)
}

// MetersToLightyears converts meters to lightyears.
func MetersToLightyears(x *big.Rat) float64 {
	ly := rat.Div(x, lyInMeters)
	lyf, _ := ly.Float64()
	return lyf
}

// VelocityToPercentC returns the percent lightspeed for a given velocity.
func VelocityToPercentC(v *big.Rat) float64 {
	pct, _ := rat.Div(v, C).Float64()
	return pct
}

// PercentCToVelocity computes the velocity in m/s from a percent of lightspeed.
func PercentCToVelocity(pct float64) *big.Rat {
	cfrac := rat.Float(pct)
	return rat.Mul(cfrac, C)
}

// AstronomicalUnit return the number of meters for a given number of AU.
func AstronomicalUnit(au float64) *big.Rat {
	return rat.Mul(rat.Float(au), oneAU)
}

// ToAstronomicalUnit converts a distance in meters to AU.
func ToAstronomicalUnit(x *big.Rat) *big.Rat {
	return rat.Div(x, oneAU)
}

// DistanceString prints out the distance in a human-readable form.
func DistanceString(d *big.Rat) string {
	switch {
	case d.Cmp(hundredthLY) > 0:
		return fmt.Sprintf("%0.4f ly", MetersToLightyears(d))
	// technically, au should be used in relation to the sun...
	case d.Cmp(tenthAU) > 0:
		return fmt.Sprintf("%s au", rat.Div(d, oneAU).FloatString(1))
	case d.Cmp(hundredthAU) > 0:
		return fmt.Sprintf("%s au", rat.Div(d, oneAU).FloatString(2))
	default:
		return rat.Div(d, rat.K).FloatString(1) + " km"
	}
}

const (
	secondsInDay  = 86400
	secondsInYear = 365.25 * secondsInDay
)

// TimeString prints out a time duration in a human-readable form. The
// Go stdlib version stops at seconds, but for printing clock drift, we'll
// need to go further.
func TimeString(dur time.Duration) string {
	var s []string
	delta := math.Round(dur.Seconds())

	years := delta / secondsInYear
	if years >= 1 {
		s = append(s, fmt.Sprintf("%0.0fy", math.Floor(years)))
		delta -= math.Floor(years) * secondsInYear
	}

	days := delta / secondsInDay
	if days >= 1 {
		s = append(s, fmt.Sprintf("%0.0fd", math.Floor(days)))
		delta -= math.Floor(days) * secondsInDay
	}

	if years < 1 && days < 7 {
		hours := delta / 3600.0
		if hours > 0 {
			s = append(s, fmt.Sprintf("%0.0fh", math.Floor(hours)))
			delta -= math.Floor(hours) * 3600.0
		}

		mins := delta / 60.0
		if mins > 0 {
			s = append(s, fmt.Sprintf("%0.0fm", math.Floor(mins)))
			delta -= math.Floor(mins) * 60.0
		}

		if delta > 0 {
			s = append(s, fmt.Sprintf("%0.0fs", math.Floor(delta)))
		}

		if len(s) == 0 {
			return "0s"
		}
	}

	return strings.Join(s, " ")
}
