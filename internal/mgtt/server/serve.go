package server

import (
	"net"
)

type Listener struct {
	address  string
	listener net.Listener
}

func Create(hostname, port string) (listener Listener) {
	listener.address = hostname + ":" + port
	return
}
