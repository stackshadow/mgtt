package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func (broker *Broker) onPacketSubscribe(connectedClient *client.MgttClient, packet *packets.SubscribePacket) (err error) {

	// PLUGINS: call CallOnSubscriptionRequest - check if subscription is accepted
	var topicResuls []byte
	for topicIndex, topic := range packet.Topics {
		qos := packet.Qoss[topicIndex]

		// call plugins
		if plugin.CallOnSubscriptionRequest(connectedClient.ID(), connectedClient.Username(), topic) == true {
			topicResuls = append(topicResuls, qos)
			connectedClient.SubScriptionAdd(topic)

			// if clean session is false, we store the subscription
			if connectedClient.CleanSessionGet() == false {
				persistance.SubscriptionsSet(
					connectedClient.ID(),
					connectedClient.Subscriptions(),
				)
			}

		} else {
			topicResuls = append(topicResuls, client.SubackErr)
		}
	}

	// thats all, respond
	connectedClient.SendSuback(packet.MessageID, topicResuls)

	// [MQTT-3.3.1-6]
	// check if an retained message exist and send it to the client
	if err == nil { // prevent multiple return
		persistance.PacketIterate("retained", func(info persistance.PacketInfo) {

			// create a packet
			publishPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
			publishPacket.MessageID = info.MessageID
			publishPacket.Retain = false
			publishPacket.Dup = info.OriginClientID != ""
			publishPacket.TopicName = info.Topic
			publishPacket.Payload = info.Payload
			publishPacket.Qos = info.Qos

			connectedClient.Publish(publishPacket)
		})
	}

	return
}
