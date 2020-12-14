package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handlePubrecPacket(connectedClient *client.MgttClient, packet *packets.PubrecPacket) (err error) {

	// we need to store the info that we get pubrec
	if qosinfo, exist := broker.pubrecs[packet.MessageID]; exist {
		log.Debug().
			Uint16("pid", packet.MessageID).
			Uint16("opid", qosinfo.originalID).
			Msg("Get pubrec and set received to true")

		qosinfo.receivedPubRec = true
		broker.pubrecs[packet.MessageID] = qosinfo

		// remove it from resend-request
		broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)

		// request pubrel
		connectedClient.SendPubrel(packet.MessageID)
	} else {
		log.Error().
			Uint16("pid", packet.MessageID).
			Msg("Not found in the pubrec-list")
	}

	return
}
