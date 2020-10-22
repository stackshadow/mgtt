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
			Uint16("mid", event.packet.Details().MessageID).
			Uint8("Qos", event.packet.Details().Qos).
			Str("packet", event.packet.String()).
			Msg("Handle packet")

		var err error = nil
		switch event.packet.(type) {

		case *packets.ConnectPacket:
			err = broker.handleConnectPacket(event)

		case *packets.SubscribePacket:
			err = broker.handleSubscribePacket(event)

		case *packets.PingreqPacket:
			err = broker.handlePingreqPacket(event)

		case *packets.PublishPacket:
			err = broker.handlePublishPacket(event)

		case *packets.PubackPacket:
			err = broker.handlePubacPacket(event)

		case *packets.PubrecPacket:
			err = broker.handlePubrecPacket(event)

		case *packets.PubrelPacket:
			err = broker.handlePubrelPacket(event)

		case *packets.PubcompPacket:
			err = broker.handlePubcompPacket(event)

		}

		if err != nil {
			log.Error().Err(err).Send()
		}

	}
}
