package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPubRel(client *client.MgttClient, packet *packets.PubrelPacket) (err error) {

	client.SendPubcomp(packet.MessageID)
	return
}
