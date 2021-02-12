package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

// PublishPacket publish a packet to all subscribers
//
// err will return the last occured error of an subscriber
func (broker *Broker) PublishPacket(packet *packets.PublishPacket, once bool) (messagedelivered bool, subscribed bool, err error) {

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
	// because it matches an established subscription
	packet.Retain = false

	messagedelivered, subscribed, err = clientlist.PublishToAllClients(packet, once)

	return
}
