//go:build !heroku

package comms

import (
	"errors"
	"log"

	"git.sr.ht/~kisom/proxima/mission"
	"gopkg.in/zeromq/goczmq.v4"
)

const (
	flagSingleMessage = 0
)

type Broadcaster struct {
	addr string
	sock *goczmq.Sock
}

func New(opts ...Opt) *Broadcaster {
	b := &Broadcaster{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Broadcaster) Connect() (err error) {
	log.Println("connecting to", b.addr)
	b.sock, err = goczmq.NewPub(b.addr)
	if err != nil {
		return err
	}
	b.sock.SetConflate(1)

	return nil
}

func (b *Broadcaster) TransmitUpdate(m *mission.Mission) error {
	if b.sock == nil {
		return errors.New("comms: broadcaster has not connected")
	}

	frame, err := m.MarshalJSON()
	if err != nil {
		return err
	}

	return b.sock.SendFrame(frame, flagSingleMessage)
}
