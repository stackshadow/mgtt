package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handlePublishPacketQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var packetID uint16 = packet.MessageID

	//  QoS1/QoS2 - Store package only if its from a real client
	if client != nil && packet.Dup == false {

		// we need a new ID
		broker.lastIDLock.Lock()

		options := messagestore.StoreResendPacketOption{
			BrokerMessageID: broker.lastID + 1,
			ClientID:        client.ID(),
			ResendAt:        time.Now().Add(time.Minute * 1),
			Packet:          packet,
		}

		//
		if packet.Qos == 1 {
			err = broker.retainedMessages.StoreResendPacket("resend", &options)
		} else {
			err = broker.retainedMessages.StoreResendPacket("unreleased", &options)
		}

		// because we stored the original message with the original messageID, we can now manipulate it
		broker.lastID = options.BrokerMessageID
		packet.MessageID = broker.lastID
		broker.lastIDLock.Unlock()
	}

	// on QoS2 we end here
	if packet.Qos == 2 {
		if client != nil {
			client.SendPubrec(packetID)
		}
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
		} else {
			broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)
		}

	}

	if client != nil {
		client.SendPuback(packetID)
	}

	return
}
