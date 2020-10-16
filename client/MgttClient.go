package client

import (
	"io"
	"net"
)

// MgttClient represents a mqtt-client
type MgttClient struct {
	id         string
	connection io.ReadWriter

	subscriptionTopics []string
}

// New create a new MgttClient with id of "unknown"
func New(connection net.Conn) (newClient *MgttClient) {

	newClient = &MgttClient{
		id:         "unknown",
		connection: connection,
	}

	return
}

// IDSet set the clientID
func (c *MgttClient) IDSet(id string) {
	c.id = id
}

// ID return the id of an MgttClient
func (c *MgttClient) ID() string {
	return c.id
}
