package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) loopHandleClientPackets() {
	for {
		event := <-broker.clientEvents

		log.Debug().
			Uint16("mid", event.packet.Details().MessageID).
			Uint8("Qos", event.packet.Details().Qos).
			Str("packet", event.packet.String()).
			Msg("Handle packet")

		var err error = nil

		switch event.packet.(type) {
		case *packets.ConnackPacket:
			err = broker.handleConackPacket(event)
		}

		switch packet := event.packet.(type) {

		case *packets.PublishPacket:
			plugin.CallOnIncoming(event.client.ID(), packet.TopicName, string(packet.Payload))

			if packet.Qos == client.SubackQoS1 {
				event.client.SendPuback(packet.MessageID)
			}

		}

		if err != nil {
			log.Error().Err(err).Send()
		}

	}
}
