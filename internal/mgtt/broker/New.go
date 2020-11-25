package broker

import (
	"gitlab.com/mgtt/internal/mgtt/client"
)

// Current is the last created broker
var Current *Broker = nil

// New will create a new Broker
func New() (broker *Broker, err error) {
	broker = &Broker{
		clients:      make(map[string]*client.MgttClient),
		clientEvents: make(chan *Event, 10),
	}

	// remember the current broker
	Current = broker

	return
}
