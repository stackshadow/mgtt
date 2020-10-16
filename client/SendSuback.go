package client

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
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
func (evt *Event) SendSuback(ReturnCodes []byte) (err error) {

	// convert
	subscr, ok := evt.Packet.(*packets.SubscribePacket)
	if ok == false {
		err = errors.New("Package is not packets.SubscribePacket")
		return
	}

	// construct the package
	suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
	suback.MessageID = subscr.MessageID
	suback.ReturnCodes = ReturnCodes

	// send it
	err = suback.Write(evt.Client.connection)

	return
}
