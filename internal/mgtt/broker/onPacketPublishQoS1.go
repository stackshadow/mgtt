package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPublishQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var ogiginalClientID string = client.ID()
	var originalPacketID uint16 = packet.MessageID
	var packetInfo persistance.PacketInfo
	var packetInfoExist bool = false

	// check if we already get this message ( e.g. LOST of PUBREC  )
	packetInfoExist, packetInfo, _ = persistance.PacketExist("qos", persistance.PacketFindOpts{
		OriginClientID:  &ogiginalClientID,
		OriginMessageID: &originalPacketID,
	})

	//  Store package ( mode A )
	if packet.Dup == false && packetInfoExist == false {

		packetInfo = persistance.PacketInfo{
			OriginClientID:  client.ID(),
			OriginMessageID: packet.MessageID,

			Topic:   packet.TopicName,
			Payload: packet.Payload,
			Qos:     packet.Qos,
		}
		err = persistance.PacketStore("qos", &packetInfo)
	}

	// Publish to all clients
	if err == nil && packetInfoExist == false {

		// create a packet
		pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
		pubPacket.MessageID = packetInfo.MessageID
		pubPacket.Retain = false
		pubPacket.Dup = packet.Dup
		pubPacket.TopicName = packetInfo.Topic
		pubPacket.Payload = packetInfo.Payload
		pubPacket.Qos = packetInfo.Qos

		// publish packet to all subscribers
		var subscribed bool
		_, subscribed, err = broker.PublishPacket(client.ID(), pubPacket, true)

		// nobody subscribe to this
		if subscribed == false {
			log.Info().
				Str("topic", pubPacket.TopicName).
				Uint16("mid", pubPacket.MessageID).
				Msg("Nobody is interested in this message")

			persistance.PacketDelete("qos", persistance.PacketFindOpts{
				OriginClientID:  &ogiginalClientID,
				OriginMessageID: &originalPacketID,
			})
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
