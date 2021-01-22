package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	messagestore "gitlab.com/mgtt/internal/mgtt/messageStore"
)

func (broker *Broker) onPacketPublishQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var originalPacketID uint16 = packet.MessageID

	//  QoS1 - Store package only if its not duplicated
	if packet.Dup == false {

		// we need a new ID
		broker.lastIDLock.Lock()

		options := messagestore.PacketInfo{
			ResendAt:  time.Now().Add(time.Minute * 1),
			Topic:     packet.TopicName,
			MessageID: broker.lastID,
			Qos:       packet.Qos,
			Payload:   packet.Payload,
		}

		//
		err = broker.retainedMessages.StoreResendPacket("resend", &options)

		broker.lastID = options.MessageID
		broker.lastIDLock.Unlock()

		// and we will send the packet with our messageID :)
		packet.MessageID = options.MessageID
	}

	// Publish to all clients
	var messagedelivered bool
	if err == nil {

		// publish packet to all subscribers
		messagedelivered, err = broker.PublishPacket(packet, false)

		// message not delivered and no error occured -> client is not interested
		if messagedelivered == false && err == nil {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")
		}

	}

	// we handle the packet, so we can ack it
	if packet.Qos == 1 {
		client.SendPuback(originalPacketID)
	}
	if packet.Qos == 2 {

		// store it in the list
		broker.pubrecs[packet.MessageID] = Qos2{
			originalClientID: client.ID(),
			originalID:       originalPacketID,
			receivedPubRec:   false,
		}
		log.Info().
			Str("topic", packet.TopicName).
			Uint16("mid", packet.MessageID).
			Msg("Store to pubrec")

		client.SendPubrec(originalPacketID)
	}

	return
}
