package client

import "github.com/eclipse/paho.mqtt.golang/packets"

func (c *MgttClient) LastWillSet(packet *packets.PublishPacket) {
	c.lastWillPacket = packet
}

func (c *MgttClient) LastWillGet() (packet *packets.PublishPacket) {
	return c.lastWillPacket
}
