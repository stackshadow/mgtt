package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// 0x00 - Success - Maximum QoS 0
// 0x01 - Success - Maximum QoS 1
// 0x02 - Success - Maximum QoS 2
// 0x80 - Failure

const (
	// SubackQoS0 ReturnCodes QoS0 for SubackPacket
	SubackQoS0 = 0x00
	// SubackQoS1 ReturnCodes QoS1 for SubackPacket
	SubackQoS1 = 0x01
	// SubackQoS2 ReturnCodes QoS2 for SubackPacket
	SubackQoS2 = 0x02
	// SubackErr ReturnCodes Err for SubackPacket
	SubackErr = 0x80
)

// SendSuback will send an SUBACK-Package
func (client *MgttClient) SendSuback(packet *packets.SubscribePacket, ReturnCodes []byte) (err error) {

	// construct the package
	suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	suback.MessageID = packet.MessageID
	suback.ReturnCodes = ReturnCodes

	log.Debug().
		Str("cid", client.ID()).
		Uint16("mid", packet.MessageID).
		Bytes("return code", ReturnCodes).
		Msg("Send PUBCOMP")

	// send it
	err = suback.Write(client.connection)

	return
}
