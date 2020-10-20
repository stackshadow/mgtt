package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPubrec will send an PUBREC-Package for QoS 2
func (client *MgttClient) SendPubcomp(MessageID uint16) (err error) {

	// construct the package
	pubcomp := packets.NewControlPacket(packets.Pubcomp).(*packets.PubcompPacket)
	pubcomp.MessageID = MessageID

	// send it
	err = pubcomp.Write(client.connection)

	return
}
