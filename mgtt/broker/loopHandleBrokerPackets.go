package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) loopHandleBrokerPackets() {
	for {
		event := <-broker.clientEvents

		log.Debug().
			Uint16("mid", event.packet.Details().MessageID).
			Uint8("Qos", event.packet.Details().Qos).
			Str("packet", event.packet.String()).
			Msg("Handle packet")

		var err error = nil

		// CONNACK-Packet
		switch event.packet.(type) {
		case *packets.ConnectPacket:
			err = broker.handleConnectPacket(event)
		}

		// check if client connects correctly
		if event.client != nil {
			if event.client.Connected == false {
				log.Error().Msg("Client not send an CONECT-Packet")
				// TODO: Force disconnect of the client
				continue
			}
		}

		switch recvPacket := event.packet.(type) {

		case *packets.SubscribePacket:
			err = broker.handleSubscribePacket(event)

		case *packets.PingreqPacket:
			err = broker.handlePingreqPacket(event)

		case *packets.PublishPacket:
			err = broker.handlePublishPacket(event.client, recvPacket)

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
