package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPubcomp(connectedClient *client.MgttClient, packet *packets.PubcompPacket) (err error) {

	log.Debug().
		Str("cid", connectedClient.ID()).
		Uint16("pid", packet.MessageID).
		Msg("We remember that we get an pubcomp")

	err = persistance.PacketPubCompSet(packet.MessageID, true)

	delete(broker.pubrecs, packet.MessageID)
	return
}
