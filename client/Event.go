package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Event represents an mgtt-event
type Event struct {
	Client *MgttClient
	Packet packets.ControlPacket
}
