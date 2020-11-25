package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handlePubacPacket(connectedClient *client.MgttClient, packet *packets.PubackPacket) (err error) {

	broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)

	return
}
