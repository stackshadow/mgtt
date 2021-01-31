package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPubRel(connectedClient *client.MgttClient, packet *packets.PubrelPacket) (err error) {
	connectedClient.SendPubcomp(packet.MessageID)
	return
}
