package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func (broker *Broker) handlePubcompPacket(event *Event) (err error) {

	// check package
	pubcompPacket, ok := event.packet.(*packets.PubcompPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	err = broker.retainedMessages.DeletePacketWithID("resend", pubcompPacket.MessageID)

	return
}
