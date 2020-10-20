package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPuback will send an PUBACK-Package
func (client *MgttClient) SendPuback(MessageID uint16) (err error) {

	// construct the package
	puback := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
	puback.MessageID = MessageID

	// send it
	err = puback.Write(client.connection)

	return
}
