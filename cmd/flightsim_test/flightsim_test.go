package flightsim_test

import (
	"log"
	"testing"
	"time"

	"git.sr.ht/~kisom/proxima/mission"
	"git.sr.ht/~kisom/proxima/physics"
	"git.sr.ht/~kisom/proxima/rat"
)

func flyMission() *mission.Mission {
	stage := mission.ActionExplore
	conn := mission.Initialize()
	start := time.Now()

	for conn.InFlight() {

		conn.Plan(time.Minute)
		if cmpStage := conn.Stage(); stage != cmpStage {
			stage = cmpStage
			now := time.Now()
			log.Println(now.Sub(start))
			log.Print(conn)
			log.Printf("Clock drift: %s", physics.TimeString(conn.Drift()))
		}
	}

	log.Println(conn)
	log.Printf("Flight time: %s",
		physics.TimeString(conn.Elapsed()))
	return conn
}

func TestMission(t *testing.T) {
	conn := flyMission()

	distance := rat.Sub(mission.ProximaDistance, conn.Distance())
	distance = physics.ToAstronomicalUnit(distance)
	if distance.Cmp(physics.AstronomicalUnit(5)) >= 0 {
		t.Logf("Mission should have ended < 5 AU from Proxima Centauri")
		t.Fatalf("Distance is %s", distance.FloatString(4))
	}
}

func BenchmarkMission(b *testing.B) {
	b.StartTimer()
	conn := flyMission()

	distance := rat.Sub(mission.ProximaDistance, conn.Distance())
	log.Println(distance)
	log.Println(">>>>", physics.DistanceString(distance))
	distance = physics.ToAstronomicalUnit(distance)
	if distance.Cmp(physics.AstronomicalUnit(5)) >= 0 {
		b.Logf("Mission should have ended < 5 AU from Proxima Centauri")
		b.Fatalf("Distance is %s", distance.FloatString(4))
	}
}
