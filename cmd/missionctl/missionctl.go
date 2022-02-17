package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"git.sr.ht/~kisom/goutils/config"
	"git.sr.ht/~kisom/proxima/handler"
)

// missionctl is the remote web server that displays updates from the mission.

func defaultAddr() string {
	host := config.GetDefault("HOST", "")
	port := config.GetDefault("PORT", "8081")
	return net.JoinHostPort(host, port)
}

func main() {
	log.Printf("PORT: '%s'", config.Get("PORT"))
	addr := flag.String("a", defaultAddr(), "`address` to listen on")
	configFile := flag.String("f", "", "`path` to config file")
	flag.Parse()

	if *configFile != "" {
		err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	srv, err := handler.NewUpstream()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", srv)
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
