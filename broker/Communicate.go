package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) Communicate() {
	for {
		event := <-broker.clientEvents

		log.Debug().
			Uint16("mid", event.Packet.Details().MessageID).
			Uint8("Qos", event.Packet.Details().Qos).
			Str("packet", event.Packet.String()).
			Msg("Received packet")

		var err error = nil
		switch event.Packet.(type) {

		case *packets.ConnectPacket:
			err = broker.handleConnectPacket(event)

		case *packets.SubscribePacket:
			err = broker.handleSubscribePacket(event)

		case *packets.PingreqPacket:
			err = broker.handlePingreqPacket(event)

		case *packets.PublishPacket:
			err = broker.handlePublishPacket(event)

		}

		if err != nil {
			log.Error().Err(err).Send()
		}

	}
}
