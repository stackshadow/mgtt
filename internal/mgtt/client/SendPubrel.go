package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPubrel will send an PUBREL-Package for QoS 2
func (client *MgttClient) SendPubrel(MessageID uint16) (err error) {

	// construct the package
	packet := packets.NewControlPacket(packets.Pubrel).(*packets.PubrelPacket)
	packet.MessageID = MessageID

	log.Debug().
		Str("cid", client.ID()).
		Uint16("mid", MessageID).
		Msg("Send PUBREL")

	// queue packet
	client.sendPackets <- packet

	return
}
