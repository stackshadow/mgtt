package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPublishQoS0(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	// Publish to all clients

	if err == nil {

		// [MQTT-3.3.1-9]
		// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
		// because it matches an established subscription
		packet.Retain = false

		// publish packet to all subscribers, except us !
		var subscribed bool
		_, subscribed, err = broker.PublishPacket(client.ID(), packet, false)

		// no message delivered
		if subscribed == false {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")
		}

	}

	return
}
