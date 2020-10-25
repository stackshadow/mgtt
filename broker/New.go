package broker

import (
	"gitlab.com/mgtt/cli"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// New will create a new Broker
func New() (broker *Broker, err error) {
	broker = &Broker{
		clients:      make(map[string]*client.MgttClient),
		clientEvents: make(chan *Event, 10),
	}

	// retainedMessages-db
	broker.retainedMessages, err = messagestore.Open(cli.CLI.DBFilename)
	if err != nil {
		return
	}

	return
}
