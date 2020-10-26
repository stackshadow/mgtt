package client

import (
	"io"
	"net"
	"time"
)

// MgttClient represents a mqtt-client
type MgttClient struct {
	id         string
	username   string
	connection io.ReadWriter
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
	connection.SetDeadline(time.Now().Add(time.Second * time.Duration(secondsTimeout)))

	return
}

// ResetTimeout will disable the timeout
func (c *MgttClient) ResetTimeout() {
	var connection net.Conn = c.connection.(net.Conn)
	connection.SetDeadline(time.Time{})
}

// IDSet set the clientID
func (c *MgttClient) IDSet(id string) {
	c.id = id
}

// ID return the id of an MgttClient
func (c *MgttClient) ID() string {
	return c.id
}
