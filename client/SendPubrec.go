package client

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// SendPubrec will send an PUBREC-Package for QoS 2
func (evt *Event) SendPubrec() (err error) {

	// convert
	publish, ok := evt.Packet.(*packets.PublishPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	// construct the package
	pubrec := packets.NewControlPacket(packets.Pubrec).(*packets.PubrecPacket)
	pubrec.MessageID = publish.MessageID

	// send it
	err = pubrec.Write(evt.Client.connection)

	return
}
