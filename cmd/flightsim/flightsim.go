package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"git.sr.ht/~kisom/proxima/handler"
	"git.sr.ht/~kisom/proxima/mission"
	"git.sr.ht/~kisom/proxima/physics"
)

func main() {
	addr := flag.String("a", "localhost:6060", "debug server `address`")
	flag.Parse()

	stage := mission.ActionExplore
	conn := mission.Initialize()
	start := time.Now()

	log.Println("starting HTTP server")
	if *addr != "" {
		srv := handler.NewStatus(conn)
		http.Handle("/", srv)
		go func() {
			log.Println(http.ListenAndServe(*addr, nil))
		}()
	}

	log.Println("underway...")

	for conn.InFlight() {

		conn.Plan(time.Minute)
		if cmpStage := conn.Stage(); stage != cmpStage {
			stage = cmpStage
			now := time.Now()
			log.Println(now.Sub(start))
			log.Println(conn)
			log.Printf("Clock drift: %s", physics.TimeString(conn.Drift()))
		}
	}

	fmt.Println(conn)
}
