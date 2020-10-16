package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleSubscribePacket(event *client.Event) {
	packet, ok := event.Packet.(*packets.SubscribePacket)
	if ok == false {
		log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
		return
	}
	log.Debug().
		Str("clientid", event.Client.ID()).
		Str("packet", packet.String()).
		Msg("RCV SubscribePacket")

	// call CallOnSubscriptionRequest - check if subscription is accepted
	var topicResuls []byte
	for topicIndex, topic := range packet.Topics {
		qos := packet.Qoss[topicIndex]

		if plugin.CallOnSubscriptionRequest(event.Client.ID(), topic) == true {
			topicResuls = append(topicResuls, qos)
			event.Client.SubScriptionAdd(topic)
		} else {
			topicResuls = append(topicResuls, client.SubackErr)
		}
	}

	// thats all, respond
	event.SendSuback(topicResuls)

	// [MQTT-3.3.1-6]
	// check if an retained message exist and send it to the client
	broker.retainedMessages.IterateRetainedTopics(func(retainedPacket *packets.PublishPacket) {
		for _, client := range broker.clients {

			// [MQTT-3.3.1-8]
			// When sending a PUBLISH Packet to a Client the Server MUST set the RETAIN flag to 1
			// if a message is sent as a result of a new subscription being made by a Client.
			retainedPacket.Retain = true

			client.Publish(retainedPacket)
		}
	})

	return
}
