package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// PacketRead reads a single packet from the connection
func (c *MgttClient) PacketRead() (packet packets.ControlPacket, err error) {
	return packets.ReadPacket(c.connection)
}
