package broker

import (
	"errors"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) handlePublishPacket(event *Event) (err error) {

	// check package
	publishPacket, ok := event.packet.(*packets.PublishPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}
	publishPacketID := publishPacket.MessageID

	if err == nil { // prevent multiple return
		// retain message ?
		if publishPacket.Retain == true && publishPacket.Dup == false {

			// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
			// [MQTT-3.3.1-11] A zero byte retained message MUST NOT be stored as a retained message on the Server.
			if len(publishPacket.Payload) == 0 {
				err = broker.retainedMessages.DeletePacketWithTopic("retained", publishPacket.TopicName)
			} else {
				// [MQTT-3.3.1-5]
				err = broker.retainedMessages.StorePacketWithTopic("retained", publishPacket.TopicName, publishPacket)
			}
		}
	}

	//  QoS1/QoS2 - Store package only if its from a real client
	if event.client != nil && (publishPacket.Qos == client.SubackQoS1 || publishPacket.Qos == client.SubackQoS2) && publishPacket.Dup == false {

		// we need a new ID
		broker.lastIDLock.Lock()

		options := messagestore.StoreResendPacketOption{
			BrokerMessageID: broker.lastID + 1,
			ClientID:        event.client.ID(),
			ResendAt:        time.Now().Add(time.Minute * 1),
			Packet:          publishPacket,
		}

		// on QoS-2, not released yet
		if publishPacket.Qos == client.SubackQoS1 {
			err = broker.retainedMessages.StoreResendPacket("resend", &options)
		} else {
			err = broker.retainedMessages.StoreResendPacket("unreleased", &options)
		}

		// because we stored the original message with the original messageID, we can now manipulate it
		broker.lastID = options.BrokerMessageID
		publishPacket.MessageID = broker.lastID
		broker.lastIDLock.Unlock()
	}

	if publishPacket.Qos == client.SubackQoS2 {
		if event.client != nil {
			event.client.SendPubrec(publishPacketID)
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
		publishPacket.Retain = false

		for _, client := range broker.clients {
			published, err = client.Publish(publishPacket)
			messagedelivered = messagedelivered || published
		}

		// no message delivered
		if messagedelivered == false {
			log.Info().
				Str("topic", publishPacket.TopicName).
				Uint16("mid", publishPacket.MessageID).
				Msg("Nobody is interested in this message")
		} else {
			broker.retainedMessages.DeletePacketWithID("resend", publishPacket.MessageID)
		}

	}

	if publishPacket.Qos == client.SubackQoS1 {
		if event.client != nil {
			event.client.SendPuback(publishPacketID)
		}
	}

	return
}
