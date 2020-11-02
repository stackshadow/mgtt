package broker

import (
	"gitlab.com/mgtt/client"
)

// New will create a new Broker
func New() (broker *Broker, err error) {
	broker = &Broker{
		clients:      make(map[string]*client.MgttClient),
		clientEvents: make(chan *Event, 10),
	}

	return
}
