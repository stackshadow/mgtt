package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPubACK(connectedClient *client.MgttClient, packet *packets.PubackPacket) (err error) {

	// TODO
	// broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)

	return
}
