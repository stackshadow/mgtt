package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPubrec will send an PUBREC-Package for QoS 2
func (client *MgttClient) SendPubrec(MessageID uint16) (err error) {

	// construct the package
	pubrec := packets.NewControlPacket(packets.Pubrec).(*packets.PubrecPacket)
	pubrec.MessageID = MessageID

	// send it
	err = pubrec.Write(client.connection)

	return
}
