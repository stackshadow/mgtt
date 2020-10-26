package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handlePublishPacketQoS0(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	// client can not be nil
	if client == nil {
		err = errors.New("Client can not be nil")
		return
	}

	// Publish to all clients
	var published bool
	var messagedelivered bool
	if err == nil {

		// [MQTT-3.3.1-9]
		// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
		// because it matches an established subscription
		packet.Retain = false

		// PLUGINS: call CallOnPublishRequest - check if publish is accepted
		for _, client := range broker.clients {
			if plugin.CallOnPublishSendRequest(client.ID(), client.Username(), packet.TopicName) == true {
				published, err = client.Publish(packet)
				messagedelivered = messagedelivered || published
			}
		}

		// no message delivered
		if messagedelivered == false {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")
		}

	}

	return
}
