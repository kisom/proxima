package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~kisom/goutils/config"
	"git.sr.ht/~kisom/proxima/comms"
	"git.sr.ht/~kisom/proxima/database"
	"git.sr.ht/~kisom/proxima/handler"
	"git.sr.ht/~kisom/proxima/mission"
)

const (
	updateInterval = time.Second
	syncInterval   = 7 * 24 * 3600 // one week
)

func updateConn(conn *mission.Mission, start time.Time) {
	ticker := time.NewTicker(updateInterval)

	for {
		select {
		case <-ticker.C:
			update := time.Now()
			conn.Plan(update.Sub(start))
			start = update
		}
	}
}

func syncClock(conn *mission.Mission) {
	ticker := time.NewTicker(syncInterval)

	for {
		select {
		case <-ticker.C:
			conn.SyncClock()
		}
	}
}

func main() {
	addr := flag.String("a", "localhost:8080",
		"`address` for status endpoint")
	configFile := flag.String("f", "/etc/proxima/proxima.conf", "`path` to config file")
	flag.Parse()

	if *configFile != "" {
		err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		db, err := database.Connect(ctx)
		if err != nil {
			log.Printf("failed to clear database: %s", err)
		}

		err = handler.ClearAllUpdates(ctx, db)
		if err != nil {
			log.Printf("failed to clear database: %s", err)
		}

		db.Close()
	}

	conn := mission.Initialize()

	// Update the spacecraft's position.
	go updateConn(conn, time.Now())

	// Sync the system clock.
	go syncClock(conn)

	broadcaster := comms.New(comms.DefaultOpts...)
	err := broadcaster.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if *addr != "" {
		srv := handler.NewStatus(conn)
		http.Handle("/", srv)
		go func() {
			log.Println(http.ListenAndServe(*addr, nil))
		}()
	}
	log.Println("underway...")

	for {
		go transmitUpdate(broadcaster, conn)
		time.Sleep(conn.DrawInterval())
	}
}

func transmitUpdate(broadcaster *comms.Broadcaster, conn *mission.Mission) {
	if err := broadcaster.TransmitUpdate(conn); err != nil {
		log.Printf("transmit failed: %s", err)
	}
}
