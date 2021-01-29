package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPubRel(connectedClient *client.MgttClient, packet *packets.PubrelPacket) (err error) {

	var pubcomp bool
	var origMessageID uint16

	// pubrel contains the original packet ID, we try to find it
	if pubcomp, origMessageID, err = persistance.PacketPubCompIsSet(packet.MessageID); err == nil {
		if pubcomp == true {
			err = connectedClient.SendPubcomp(origMessageID)
		} else {
			err = errors.New("Not received an pubcomp for this packet")
		}
	}

	if err == nil {
		connectedClient.SendPubcomp(origMessageID)

		// remove the message from the store
		
	} else {
		log.Error().
			Err(err).
			Str("cid", connectedClient.ID()).
			Uint16("pid", packet.MessageID).
			Send()
	}

	return
}
