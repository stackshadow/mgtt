package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePingreqPacket(event *client.Event) (err error) {

	// check package
	_, ok := event.Packet.(*packets.PingreqPacket)
	if ok == false {
		err = errors.New("Package is not packets.PingreqPacket")
		return
	}

	err = event.Client.SendPingresp()
	return
}
