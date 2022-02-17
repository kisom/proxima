package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"git.sr.ht/~kisom/proxima/mission"
)

// Status runs an HTTP endpoint that gives status updates.
type Status struct {
	m *mission.Mission
}

func NewStatus(m *mission.Mission) *Status {
	return &Status{
		m: m,
	}
}

// ServeHTTP sends current flight information.
func (s *Status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{
		"mission": s.m,
		"drift":   s.m.Drift().Seconds(),
		"elapsed": s.m.Elapsed().Seconds(),
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshal data to JSON: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
