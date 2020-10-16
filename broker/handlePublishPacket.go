package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePublishPacket(event *client.Event) (err error) {
	packet, ok := event.Packet.(*packets.PublishPacket)
	if ok == false {
		log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
		return
	}

	// retain message ?
	if packet.Retain == true {

		// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
		// [MQTT-3.3.1-11]  A zero byte retained message MUST NOT be stored as a retained message on the Server.
		if len(packet.Payload) == 0 {
			broker.retainedMessages.DeleteRetainedIfExist(packet.TopicName)
		} else {

			// [MQTT-3.3.1-5]
			broker.retainedMessages.StoreRetainedTopic(packet)
		}
	}

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a
	// Client because it matches an established subscription

	packet.Retain = false
	for _, client := range broker.clients {
		client.Publish(packet)
	}

	return
}
