package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendUnsubAck will send an UNSUBACK-Package
func (client *MgttClient) SendUnsubAck(MessageID uint16) (err error) {

	// construct the package
	packet := packets.NewControlPacket(packets.Unsuback).(*packets.UnsubackPacket)
	packet.MessageID = MessageID

	log.Debug().
		Str("cid", client.ID()).
		Uint16("mid", packet.MessageID).
		Msg("Send UNSUBACK")

	// queue packet
	client.sendPackets <- packet

	return
}
