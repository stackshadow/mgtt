package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPublishQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var originalPacketID uint16 = packet.MessageID

	//  QoS1 - Store package only if its not duplicated
	if packet.Dup == false {

		// we need a new ID
		broker.lastIDLock.Lock()

		err = persistance.PacketStore(
			persistance.PacketInfo{
				OriginClientID: client.ID(),
				ResendAt:       time.Now().Add(time.Minute * 1),

				Topic:     packet.TopicName,
				MessageID: packet.MessageID,
				Payload:   packet.Payload,
				Qos:       packet.Qos,
			},
			&broker.lastID,
		)

		broker.lastIDLock.Unlock()

		// and we will send the packet with our messageID :)
		packet.MessageID = broker.lastID
	}

	// Publish to all clients
	var messagedelivered bool
	if err == nil {

		// publish packet to all subscribers
		messagedelivered, err = broker.PublishPacket(packet, false)

		// message not delivered and no error occured -> no client is interested
		if messagedelivered == false && err == nil {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")

			// on QOS2 we emulate that we get PUBCOMP
			if packet.Qos == 2 {
				persistance.PacketPubCompSet(originalPacketID, true)
			}
		}

	}

	// we handle the packet, so we can ack it
	if packet.Qos == 1 {
		client.SendPuback(originalPacketID)
	}
	if packet.Qos == 2 {
		client.SendPubrec(originalPacketID)
	}

	return
}
