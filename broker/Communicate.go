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
func (broker *Broker) Communicate() {
	for {
		event := <-broker.clientEvents
		switch event.Packet.(type) {

		case *packets.ConnectPacket:
			packet, ok := event.Packet.(*packets.ConnectPacket)
			if ok == false {
				log.Error().Str("clientid", event.Client.ID()).Msg("Expected ConnectPacket")
				break
			}
			log.Debug().
				Str("clientid", event.Client.ID()).
				Msg("RCV ConnectPacket")

			// call CallOnAcceptNewClient - check if we accept the client
			accepted := plugin.CallOnAcceptNewClient(event.Client.ID(), packet.Username, string(packet.Password))
			if accepted == false {
				log.Error().Str("clientid", event.Client.ID()).Msg("Client not accepted by plugin")
				event.SendConnack(client.ConnackUnauthorized)
				break
			}

			log.Info().Str("clientid", event.Client.ID()).Msg("Add new client to client-list")
			broker.clients[event.Client.ID()] = event.Client
			event.SendConnack(client.ConnackAccepted)

		case *packets.SubscribePacket:
			packet, ok := event.Packet.(*packets.SubscribePacket)
			if ok == false {
				log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
				break
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
					client.Publish(retainedPacket)
				}
			})

		case *packets.PingreqPacket:
			packet, ok := event.Packet.(*packets.PingreqPacket)
			if ok == false {
				log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
				break
			}
			log.Debug().
				Str("clientid", event.Client.ID()).
				Str("packet", packet.String()).
				Msg("RCV PingreqPacket")

			event.SendPingresp()

		case *packets.PublishPacket:
			packet, ok := event.Packet.(*packets.PublishPacket)
			if ok == false {
				log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
				break
			}

			// retain message ?
			if packet.Retain == true {

				// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
				if len(packet.Payload) == 0 {
					broker.retainedMessages.DeleteRetainedIfExist(packet.TopicName)
				} else {

					// [MQTT-3.3.1-5]
					broker.retainedMessages.StoreRetainedTopic(packet)
				}

				// [MQTT-3.3.1-9]
				// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a
				// Client because it matches an established subscription
				packet.Retain = false
			}

			for _, client := range broker.clients {
				client.Publish(packet)
			}

		}

	}
}
