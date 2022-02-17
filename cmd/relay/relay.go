package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"git.sr.ht/~kisom/goutils/config"
	"git.sr.ht/~kisom/proxima/handler"
)

// Relay retransmits updates from the mission to an upstream web server.

func getUpdateFromHTTP(addr string) (*handler.Update, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	u, err := handler.UpdateFromReader(resp.Body)
	return u, err
}

func postUpdate(addr string, stats *handler.Update) error {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(stats)
	if err != nil {
		return err
	}

	resp, err := handler.WithBasicAuth(addr, "text/json", buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP POST failed: %d %s",
			resp.StatusCode, resp.Status)
		return fmt.Errorf("relay: HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	resp.Body.Close()
	return nil
}

func relay(local, remote string) error {
	stats, err := getUpdateFromHTTP(local)
	if err != nil {
		return err
	}

	return postUpdate(remote, stats)
}

func main() {
	addr := flag.String("d", "http://localhost:8080", "`address` of flight control")
	configFile := flag.String("f", "/etc/proxima/proxima.conf", "`path` to config file")
	upstream := flag.String("u", "http://localhost:8081", "`address` of update server")
	flag.Parse()

	if *configFile != "" {
		err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := relay(*addr, *upstream)
	if err != nil {
		log.Fatal(err)
	}
}
