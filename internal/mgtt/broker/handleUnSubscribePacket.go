package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handleUnSubscribePacket(connectedClient *client.MgttClient, packet *packets.UnsubscribePacket) (err error) {

	// unsubscribe
	connectedClient.SubScriptionsRemove(packet.Topics)

	// thats all, respond
	err = connectedClient.SendUnsubAck(packet.MessageID)

	return
}
