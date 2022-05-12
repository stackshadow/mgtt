package client

import (
	"net"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// MgttClient represents a mqtt-client
type MgttClient struct {
	id           string
	username     string
	cleanSession bool
	connection   net.Conn
	Connected    bool

	// The last will of this client
	lastWillPacket *packets.PublishPacket

	subscriptionTopics []string

	sendPackets chan packets.ControlPacket // send-buffer to avoid double-write

	// loop signals
	packetSendLoopRunning bool // indicates i the send-loop is running
	packetSendLoopExit    chan bool
}

// Init create a new MgttClient with id of "unknown"
func (c *MgttClient) Init(connection net.Conn, timeout time.Duration) {

	// create a new client with an new random-id
	guid := xid.New()

	c.id = guid.String()
	c.connection = connection
	c.sendPackets = make(chan packets.ControlPacket, 10)
	c.packetSendLoopExit = make(chan bool)

	// setup timeout
	if timeout > 0 {
		log.Debug().Dur("timeout", timeout).Msg("Set deadline for client")
		connection.SetDeadline(time.Now().Add(timeout))
	}

	// start the write-loop
	go c.packetSendLoop()

	return
}

// ResetTimeout will disable the timeout
func (c *MgttClient) ResetTimeout() {

	c.connection.SetDeadline(time.Time{})
}

// IDSet set the clientID
func (c *MgttClient) IDSet(id string) {
	c.id = id
}

// ID return the id of an MgttClient
func (c *MgttClient) ID() string {
	return c.id
}
