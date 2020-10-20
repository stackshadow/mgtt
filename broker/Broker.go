package broker

import (
	"sync"

	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// Broker represents a broker
type Broker struct {
	clients map[string]*client.MgttClient

	// clientEvents are raw incoming events qos=0 qos=1 qos=2
	clientEvents chan *Event

	// lastID of our messages
	lastID     uint16
	lastIDLock sync.Mutex

	// remember pubrecs
	// map of broker-Message-ID
	pubrec map[uint16]*messagestore.StoreResendPacketOption

	// an boltDB to store retained messages
	retainedMessages *messagestore.Store
}
