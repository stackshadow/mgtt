package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPubrec will send an PUBREC-Package for QoS 2
func (client *MgttClient) SendPubcomp(MessageID uint16) (err error) {

	// construct the package
	packet := packets.NewControlPacket(packets.Pubcomp).(*packets.PubcompPacket)
	packet.MessageID = MessageID

	log.Debug().
		Str("cid", client.ID()).
		Uint16("mid", MessageID).
		Msg("Send PUBCOMP")

	// queue packet
	client.sendPackets <- packet

	return
}
