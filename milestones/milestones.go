package milestones

import (
	"math/big"

	"git.sr.ht/~kisom/proxima/mission"
	"git.sr.ht/~kisom/proxima/rat"
)

type state struct {
	d  *big.Rat
	et float64
	dt float64
}

type check func(state) bool

func checkDistance(target *big.Rat) check {
	return func(s state) bool {
		return s.d.Cmp(target) > 0
	}
}

func checkElapsed(target float64) check {
	return func(s state) bool {
		return s.et > target
	}
}

func checkDrift(target float64) check {
	return func(s state) bool {
		return s.dt > target
	}
}

type Milestone struct {
	m string
	c check
}

func (m Milestone) String() string {
	return m.m
}

func (m Milestone) Cmp(state state) bool {
	return m.c(state)
}

var oneHundredYears = 100 * 365.25 * 86400

var milestones = []Milestone{
	{"passed the orbit of Mars", checkDistance(mission.MarsDistance)},
	{"passed the orbit of Jupiter", checkDistance(mission.JupiterDistance)},
	{"passed the orbit of Saturn", checkDistance(mission.SaturnDistance)},
	{"passed the orbit of Uranus", checkDistance(mission.UranusDistance)},
	{"passed the orbit of Neptune", checkDistance(mission.NeptuneDistance)},
	{"passed the orbit of Pluto", checkDistance(mission.PlutoDistance)},
	{"passed the termination shock", checkDistance(mission.TerminationShock)},
	{"left the solar system", checkDistance(mission.Heliopause)},
	{"everyone you know is probably dead", checkDrift(oneHundredYears)},
}

func newState(distance string, elapsed, drift float64) (s state, err error) {
	s.d = rat.FromString(distance)
	s.et = elapsed
	s.dt = drift
	return s, nil
}

// Get the relevant milestones for the current mission.
func Get(distance string, elapsed, drift float64) ([]string, error) {
	var ms []string
	s, err := newState(distance, elapsed, drift)
	if err != nil {
		return nil, err
	}

	for _, m := range milestones {
		if m.Cmp(s) {
			ms = append(ms, m.String())
		}
	}

	msl := len(ms)
	for i := 0; i < msl/2; i++ {
		j := (msl - i - 1)
		ms[i], ms[j] = ms[j], ms[i]
	}
	return ms, err
}
