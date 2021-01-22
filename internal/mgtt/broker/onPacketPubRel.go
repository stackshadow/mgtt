package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPubRel(connectedClient *client.MgttClient, packet *packets.PubrelPacket) (err error) {

	// pubrel contains the original packet ID, we try to find it
	for _, pubrec := range broker.pubrecs {
		if pubrec.originalClientID == connectedClient.ID() {
			if pubrec.originalID == packet.MessageID {
				log.Debug().
					Uint16("pid", packet.MessageID).
					Uint16("opid", pubrec.originalID).
					Msg("Found packet id in pubrec-list")

				connectedClient.SendPubcomp(packet.MessageID)
				return
			}
		}
	}

	log.Error().
		Uint16("pid", packet.MessageID).
		Msg("Not found in the pubrec-list")

	return
}
