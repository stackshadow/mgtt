package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// readPacket reads a single packet from the connection and store it to the buffer
func (c *MgttClient) readPacket() (packet packets.ControlPacket, err error) {
	return packets.ReadPacket(c.connection)
}
