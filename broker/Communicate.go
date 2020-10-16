package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) Communicate() {
	for {
		event := <-broker.clientEvents
		switch event.Packet.(type) {

		case *packets.ConnectPacket:
			broker.handleConnectPacket(event)

		case *packets.SubscribePacket:
			broker.handleSubscribePacket(event)

		case *packets.PingreqPacket:
			broker.handlePingreqPacket(event)

		case *packets.PublishPacket:
			broker.handlePublishPacket(event)

		}

	}
}
