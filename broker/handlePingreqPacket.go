package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func (broker *Broker) handlePingreqPacket(event *Event) (err error) {

	// check package
	_, ok := event.packet.(*packets.PingreqPacket)
	if ok == false {
		err = errors.New("Package is not packets.PingreqPacket")
		return
	}

	err = event.client.SendPingresp()
	return
}
