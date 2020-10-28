package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePubcompPacket(connectedClient *client.MgttClient, packet *packets.PubcompPacket) (err error) {

	err = broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)

	return
}
