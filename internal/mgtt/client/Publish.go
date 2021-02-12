package client

import (
	"strings"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// Publish an publish-packet
//
// this function check if the topic in the `packet` matches the client-topic-filter
//
// - return published=true if the message could be send
//
// - return subscribed=true if an sibscription match
//
// - err is returned if something is wrong with the connection
func (client *MgttClient) Publish(packet *packets.PublishPacket) (published bool, subscribed bool, err error) {

	topic := packet.TopicName
	topicArray := strings.Split(topic, "/")
	topicMatched := false

	// check if one of our subscription matched
	// this prevents multiple sending of packet to a single client
	for _, subscriptionTopic := range client.subscriptionTopics {
		subscriptionTopicArray := strings.Split(subscriptionTopic, "/")

		// [MQTT-3.3.2-3]
		// The Topic Name in a PUBLISH Packet sent by a Server to a subscribing Client
		// MUST match the Subscription’s Topic Filter
		if Match(subscriptionTopicArray, topicArray) == true {
			subscribed = true
			topicMatched = true
			break
		}

	}

	if topicMatched == true {
		log.Info().
			Str("clientid", client.ID()).
			Uint16("mid", packet.MessageID).
			Str("topic", packet.TopicName).
			Msg("Publish message to client")

		client.sendPackets <- packet

		published = true
	}

	return
}
