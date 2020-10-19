package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePublishPacket(event *client.Event) (err error) {

	// check package
	packet, ok := event.Packet.(*packets.PublishPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	if err == nil { // prevent multiple return
		// retain message ?
		if packet.Retain == true {

			// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
			// [MQTT-3.3.1-11] A zero byte retained message MUST NOT be stored as a retained message on the Server.
			if len(packet.Payload) == 0 {
				err = broker.retainedMessages.DeletePacketWithTopic("retained", packet.TopicName)
			} else {
				// [MQTT-3.3.1-5]
				err = broker.retainedMessages.StorePacketWithTopic("retained", packet.TopicName, packet)
			}
		}
	}

	// Handle QoS-1/2 - Reserve ID
	if packet.Qos == client.SubackQoS1 || packet.Qos == client.SubackQoS2 {
		// we need a new ID
		broker.lastID++
		broker.lastID, err = broker.retainedMessages.StorePacketWithID("resend", broker.lastID, nil)
		packet.MessageID = broker.lastID
	}

	// Publish to all clients
	//
	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a
	// Client because it matches an established subscription
	var published bool
	var messagedelivered bool = true
	if err == nil { // prevent multiple return
		packet.Retain = false

		for _, client := range broker.clients {
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

	// Handle QoS-1/2
	if packet.Qos == client.SubackQoS1 || packet.Qos == client.SubackQoS2 {

		// store packet
		if messagedelivered == true {
			broker.lastID, err = broker.retainedMessages.StorePacketWithID("resend", packet.MessageID, packet)
		} else {
			broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)
		}

		if packet.Qos == client.SubackQoS1 {
			// we ignore the returned err by purpose
			event.SendPuback()
		}
	}

	return
}
