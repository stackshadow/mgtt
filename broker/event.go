package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

// Event represents an mgtt-event
type Event struct {
	client *client.MgttClient
	packet packets.ControlPacket
}
