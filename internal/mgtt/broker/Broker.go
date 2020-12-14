package broker

import (
	"sync"

	"gitlab.com/mgtt/internal/mgtt/client"
	messagestore "gitlab.com/mgtt/internal/mgtt/messageStore"
)

var ConnectTimeout int64 = 30

// Broker represents a broker
type Broker struct {
	clients map[string]*client.MgttClient

	// clientEvents are raw incoming events qos=0 qos=1 qos=2
	clientEvents chan *Event

	// lastID of our messages
	lastID     uint16
	lastIDLock sync.Mutex

	// an boltDB to store retained messages
	retainedMessages *messagestore.Store

	// myid -> originalID
	pubrecs map[uint16]Qos2
}

// Qos2 store infos for QoS2-Packets
type Qos2 struct {
	originalClientID string
	originalID       uint16
	receivedPubRec   bool
}
