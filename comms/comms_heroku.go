//go:build heroku

package comms

import (
	"log"

	"git.sr.ht/~kisom/proxima/mission"
)

type Broadcaster struct {
	addr string
}

func New(opts ...Opt) *Broadcaster {
	b := &Broadcaster{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Broadcaster) Connect() (err error) {
	log.Println("stub broadcaster in use")
	return nil
}

func (b *Broadcaster) TransmitUpdate(m *mission.Mission) error {
	log.Println("stub broadcaster in use")
	return nil
}
