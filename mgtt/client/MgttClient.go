package client

import (
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

// MgttClient represents a mqtt-client
type MgttClient struct {
	id         string
	username   string
	connection net.Conn
	Connected  bool

	subscriptionTopics []string
}

// New create a new MgttClient with id of "unknown"
func New(connection net.Conn, secondsTimeout int64) (newClient *MgttClient) {

	newClient = &MgttClient{
		id:         "unknown",
		connection: connection,
	}

	// setup timeout
	if secondsTimeout > 0 {
		log.Debug().Int64("timeout", secondsTimeout).Msg("Set deadline for client")
		connection.SetDeadline(time.Now().Add(time.Second * time.Duration(secondsTimeout)))
	}

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
