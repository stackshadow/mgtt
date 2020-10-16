package client

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPuback will send an PUBACK-Package
func (evt *Event) SendPuback() (err error) {

	// convert
	publish, ok := evt.Packet.(*packets.PublishPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	// construct the package
	puback := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
	puback.MessageID = publish.MessageID

	// send it
	err = puback.Write(evt.Client.connection)

	return
}
