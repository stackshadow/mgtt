package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handlePubcompPacket(connectedClient *client.MgttClient, packet *packets.PubcompPacket) (err error) {
	log.Debug().
		Uint16("pid", packet.MessageID).
		Msg("Remove from pubrec-list")

	delete(broker.pubrecs, packet.MessageID)
	return
}
