package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPubrel will send an PUBREL-Package for QoS 2
func (client *MgttClient) SendPubrel(MessageID uint16) (err error) {

	// construct the package
	pubrel := packets.NewControlPacket(packets.Pubrel).(*packets.PubrelPacket)
	pubrel.MessageID = MessageID

	// send it
	err = pubrel.Write(client.connection)

	return
}
