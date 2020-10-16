package client

import (
	"strings"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

func match(route []string, topic []string) bool {
	if len(route) == 0 {
		return len(topic) == 0
	}

	if len(topic) == 0 {
		return route[0] == "#"
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return match(route[1:], topic[1:])
	}
	return false
}

func (c *MgttClient) Publish(pubpacket *packets.PublishPacket) (published bool, err error) {

	topic := pubpacket.TopicName
	topicArray := strings.Split(topic, "/")

	for _, subscriptionTopic := range c.subscriptionTopics {
		subscriptionTopicArray := strings.Split(subscriptionTopic, "/")

		// [MQTT-3.3.2-3]
		// The Topic Name in a PUBLISH Packet sent by a Server to a subscribing Client
		// MUST match the Subscription’s Topic Filter
		if match(subscriptionTopicArray, topicArray) == true {
			log.Info().
				Str("clientid", c.ID()).
				Uint16("mid", pubpacket.MessageID).
				Str("topic", pubpacket.TopicName).
				Msg("Publish message to client")

			err = pubpacket.Write(c.connection)
			published = true
			return
		}

	}

	return
}
