package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleSubscribePacket(connectedClient *client.MgttClient, packet *packets.SubscribePacket) (err error) {

	// PLUGINS: call CallOnSubscriptionRequest - check if subscription is accepted
	var topicResuls []byte
	for topicIndex, topic := range packet.Topics {
		qos := packet.Qoss[topicIndex]

		// call plugins
		if plugin.CallOnSubscriptionRequest(connectedClient.ID(), connectedClient.Username(), topic) == true {
			topicResuls = append(topicResuls, qos)
			connectedClient.SubScriptionAdd(topic)
		} else {
			topicResuls = append(topicResuls, client.SubackErr)
		}
	}

	// thats all, respond
	connectedClient.SendSuback(packet, topicResuls)

	// [MQTT-3.3.1-6]
	// check if an retained message exist and send it to the client
	if err == nil { // prevent multiple return
		broker.retainedMessages.IteratePackets("retained", func(retainedPacket *packets.PublishPacket) {
			for _, client := range broker.clients {

				// [MQTT-3.3.1-8]
				// When sending a PUBLISH Packet to a Client the Server MUST set the RETAIN flag to 1
				// if a message is sent as a result of a new subscription being made by a Client.
				retainedPacket.Retain = true

				client.Publish(retainedPacket)
			}
		})
	}

	return
}
