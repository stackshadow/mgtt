package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPubRec(connectedClient *client.MgttClient, packet *packets.PubrecPacket) (err error) {

	log.Debug().
		Str("cid", connectedClient.ID()).
		Uint16("pid", packet.MessageID).
		Msg("Get pubrec, reply with pubrel")

	// request pubrel
	connectedClient.SendPubrel(packet.MessageID)

	return
}
