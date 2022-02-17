package handler

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"time"

	"git.sr.ht/~kisom/goutils/config"
	"git.sr.ht/~kisom/proxima/milestones"
	"git.sr.ht/~kisom/proxima/mission"
	"git.sr.ht/~kisom/proxima/physics"
	"git.sr.ht/~kisom/proxima/rat"
)

const timeFormat = "2006-01-02 15:04 MST"

//go:embed templates/index.html
var indexTemplatePage string
var indexTemplate = template.Must(template.New("index").Parse(indexTemplatePage))

type clock struct {
	Observer string
	Relative string
}

type state struct {
	Velocity   string
	VelocityPC string
	Distance   string
	Remaining  string
}

type Page struct {
	NoUpdates  bool
	Timestamp  string
	Phase      string
	Elapsed    string
	Drift      string
	Clock      clock
	State      state
	Milestones []string
	Simulation bool
}

func (p *Page) HasMilestones() bool {
	return len(p.Milestones) > 0
}

func PageFromUpdate(created int64, update *Update) *Page {
	if update == nil {
		return &Page{
			NoUpdates:  true,
			Simulation: config.Get("SIMULATION") == "true",
		}
	}

	page := &Page{
		Timestamp: time.Unix(created, 0).Format(timeFormat),
		Phase:     update.Mission.Action,
		Clock: clock{
			Observer: update.Mission.Clock.Observer,
			Relative: update.Mission.Clock.Observer,
		},
	}

	page.Elapsed = physics.TimeString(rat.DurationSeconds(rat.Float(update.Elapsed)))
	page.Drift = physics.TimeString(rat.DurationSeconds(rat.Float(update.Drift)))
	v := rat.FromString(update.Mission.State.V)
	x := rat.FromString(update.Mission.State.X)

	page.State.Velocity = rat.Div(v, rat.K).FloatString(0)
	page.State.VelocityPC = fmt.Sprintf("%0.3f", physics.VelocityToPercentC(v))
	page.State.Distance = physics.DistanceString(x)
	page.State.Remaining = physics.DistanceString(rat.Sub(mission.ProximaDistance, x))
	ms, err := milestones.Get(update.Mission.State.V, update.Elapsed, update.Drift)
	if err != nil {
		log.Println(err)
		return page
	}
	page.Milestones = ms
	page.Simulation = config.Get("SIMULATION") == "true"

	return page
}
