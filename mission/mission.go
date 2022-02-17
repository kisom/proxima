package mission

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"git.sr.ht/~kisom/proxima/physics"
	"git.sr.ht/~kisom/proxima/rat"
)

const (
	proximaLY = 4.247
)

var (
	Heliopause       = physics.AstronomicalUnit(120.0)
	MarsDistance     = physics.AstronomicalUnit(0.52)
	JupiterDistance  = physics.AstronomicalUnit(4.20)
	SaturnDistance   = physics.AstronomicalUnit(8.58)
	UranusDistance   = physics.AstronomicalUnit(18.20)
	NeptuneDistance  = physics.AstronomicalUnit(29.05)
	PlutoDistance    = physics.AstronomicalUnit(38.48)
	TerminationShock = physics.AstronomicalUnit(90.0)
	fiftyAU          = physics.AstronomicalUnit(50.0)
)

var conn *Mission

type Action uint8

func (a Action) String() string {
	switch a {
	case ActionAccelerate:
		return "accelerating"
	case ActionCoast:
		return "coasting"
	case ActionDecelerate:
		return "decelerating"
	case ActionExplore:
		return "exploring"
	default:
		return "mission error"
	}
}

const (
	ActionAccelerate Action = iota + 1
	ActionCoast
	ActionDecelerate
	ActionExplore
)

const clockFormat = "2006-01-02 15:04 MST"

var (
	// ProximaDistance is the distance to proxima in meters.
	ProximaDistance            = physics.LightyearsToMeters(proximaLY)
	VelocityEscape             = rat.UInt64(11180)
	VelocityCruise             = physics.PercentCToVelocity(0.999)
	VelocityExplore            = rat.Float(18975.1)
	DecelerationTargetDistance *big.Rat
)

func decelerationTargetDistance() {
	if DecelerationTargetDistance != nil {
		return
	}

	// How long does it take to transition from cruising velocity to
	// exploration velocity?
	t := rat.Div(rat.Sub(VelocityExplore, VelocityCruise), physics.NegGee)

	// How much distance is covered during the deceleration, assuming we want
	// to end up about 5 AU from Proxima Centauri?
	dx := physics.AccelerationDistance(physics.AstronomicalUnit(5), VelocityCruise, physics.NegGee, t)

	// Compute the point at which we need to start decelerating.
	DecelerationTargetDistance = rat.Sub(ProximaDistance, dx)
}

type Mission struct {
	state  physics.State
	clock  *physics.Clock
	action Action
}

func (m *Mission) Stage() Action {
	return m.action
}

func (m *Mission) InFlight() bool {
	return m.action != ActionExplore
}

func (m *Mission) distanceFromProximaCentauri() *big.Rat {
	return rat.Sub(ProximaDistance, m.state.X)
}

func (m *Mission) String() string {
	remaining := m.distanceFromProximaCentauri()
	return fmt.Sprintf(`Phase: %s
 Ship time: %s
Earth time: %s
  Velocity: %s km/s (%0.3fc)
 Traveled: %s
Remaining: %s
`,
		m.action,
		m.clock.Observer.Now().Format(clockFormat),
		m.clock.Relative.Now().Format(clockFormat),
		rat.Div(m.state.V, rat.K).FloatString(1),
		physics.VelocityToPercentC(m.state.V),
		physics.DistanceString(m.state.X),
		physics.DistanceString(remaining),
	)
}

func (m *Mission) Lines() []string {
	remaining := m.distanceFromProximaCentauri()
	lines := make([]string, 6)
	lines[0] = "Phase: " + m.action.String()
	lines[1] = fmt.Sprintf(" Ship time: %s", m.clock.Observer.Now().Format(clockFormat))
	lines[2] = fmt.Sprintf("Earth time: %s", m.clock.Relative.Now().Format(clockFormat))
	lines[3] = fmt.Sprintf("  Velocity: %s km/s (%0.3fc)",
		rat.Div(m.state.V, rat.K).FloatString(1),
		physics.VelocityToPercentC(m.state.V))
	lines[4] = fmt.Sprintf(" Traveled: %s", physics.DistanceString(m.state.X))
	lines[5] = fmt.Sprintf("Remaining: %s", physics.DistanceString(remaining))

	return lines
}

func Initialize() *Mission {
	decelerationTargetDistance()
	conn = &Mission{
		state: physics.State{
			X: rat.Zero,
			V: VelocityEscape,
		},
		clock:  physics.Now(),
		action: ActionAccelerate,
	}
	return conn
}

func (m *Mission) Plan(d time.Duration) {
	seconds := d.Seconds()
	switch m.action {
	case ActionAccelerate:
		m.state = m.state.Accelerate(physics.Gee, seconds)
		if m.state.V.Cmp(VelocityCruise) >= 0 {
			m.action = ActionCoast
		}
	case ActionCoast:
		m.state = m.state.Coast(seconds)
		if m.state.X.Cmp(DecelerationTargetDistance) >= 0 {
			m.action = ActionDecelerate
		}
	case ActionDecelerate:
		m.state = m.state.Accelerate(physics.NegGee, seconds)
		if m.state.V.Cmp(VelocityExplore) <= 0 {
			m.action = ActionExplore
		}
	case ActionExplore:
		// nothing to do
	}
	m.clock.Update(seconds, m.state.V)
}

func (m *Mission) DrawInterval() time.Duration {
	switch m.action {
	case ActionAccelerate:
		if m.state.X.Cmp(JupiterDistance) < 0 {
			return time.Second
		}

		if m.state.X.Cmp(Heliopause) < 0 {
			return time.Minute
		}
	case ActionDecelerate:
		remaining := rat.Sub(ProximaDistance, m.state.X)
		if remaining.Cmp(fiftyAU) < 0 {
			return time.Minute
		}

		if remaining.Cmp(JupiterDistance) < 0 {
			return time.Second
		}
	}

	return time.Hour
}

// SyncClock syncs the clock to the system time.
func (m *Mission) SyncClock() {
	m.clock.Sync()
}

func (m *Mission) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{}
	v["action"] = m.action.String()
	v["state"] = map[string]string{
		"x": m.state.X.FloatString(0),
		"v": m.state.V.FloatString(0),
	}
	v["clock"] = map[string]interface{}{
		"launched":    m.clock.Launch.Format(clockFormat),
		"observer":    m.clock.Observer.Now().Format(clockFormat),
		"observer_et": m.clock.Observer.Now().Sub(m.clock.Launch).Seconds(),
		"relative":    m.clock.Relative.Now().Format(clockFormat),
		"relative_et": m.clock.Observer.Now().Sub(m.clock.Launch).Seconds(),
	}

	return json.Marshal(v)
}

// Distance returns the mission's current distance.
func (m *Mission) Distance() *big.Rat {
	x := rat.Rat()
	x.Set(m.state.X)
	return x
}

func (m *Mission) Drift() time.Duration {
	return m.clock.Drift()
}

func (m *Mission) Elapsed() time.Duration {
	return m.clock.Observer.Since(m.clock.Launch)
}
