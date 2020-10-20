package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func (broker *Broker) handlePubacPacket(event *Event) (err error) {

	// check package
	packet, ok := event.packet.(*packets.PubackPacket)
	if ok == false {
		err = errors.New("Package is not packets.PubackPacket")
		return
	}

	broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)

	return
}
