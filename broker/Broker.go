package broker

import (
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// Broker represents a broker
type Broker struct {
	clients map[string]*client.MgttClient

	// clientEvents are raw incoming events qos=0 qos=1 qos=2
	clientEvents chan *client.Event

	// lastID of our messages
	lastID uint16
	// pending events are messages with qos=1 qos=2
	pendingEvents map[uint16]*client.Event

	// an boltDB to store retained messages
	retainedMessages *messagestore.Store
}
