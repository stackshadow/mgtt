package client

import (
	"io"
	"net"

	"github.com/google/uuid"
)

// MgttClient represents a mqtt-client
type MgttClient struct {
	id         string
	connection io.ReadWriter

	subscriptionTopics []string
}

// New create a new MgttClient with an new UUID
func New(connection net.Conn) (newClient *MgttClient) {
	// set a new UUID for the client
	newUUID := uuid.New()

	newClient = &MgttClient{
		id:         newUUID.String(),
		connection: connection,
	}

	return
}

// ID return the id of an MgttClient
func (c *MgttClient) ID() string {
	return c.id
}
