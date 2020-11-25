package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// GetPacket return a packet from the internal buffer
func (c *MgttClient) GetPacket() (packet packets.ControlPacket) {
	return <-c.recvPackets
}
