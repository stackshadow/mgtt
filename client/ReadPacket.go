package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// ReadPacket reads a single packet from the connection
//
// return err if an error occurred
func (c *MgttClient) ReadPacket() (packet packets.ControlPacket, err error) {
	return packets.ReadPacket(c.connection)
}
