package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPuback will send an PUBACK-Package
func (client *MgttClient) SendPuback(MessageID uint16) (err error) {

	// construct the package
	puback := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
	puback.MessageID = MessageID

	log.Debug().
		Str("cid", client.ID()).
		Uint16("mid", MessageID).
		Msg("Send PUBACK")

	// send it
	err = puback.Write(client.connection)

	return
}
