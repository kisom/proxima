package comms

import (
	"fmt"
	"net"
)

type Opt func(b *Broadcaster)

var DefaultOpts = []Opt{
	TCPAddress("127.0.0.1", "4247"),
}

// TCPAddress applies the host and port to the broadcaster.
func TCPAddress(host, port string) Opt {
	addr := fmt.Sprintf("tcp://%s", net.JoinHostPort(host, port))
	return func(b *Broadcaster) {
		b.addr = addr
	}
}
