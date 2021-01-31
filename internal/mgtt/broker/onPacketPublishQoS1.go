package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPublishQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var ogiginalPacketID string = client.ID()
	var originalPacketID uint16 = packet.MessageID
	var packetInfo persistance.PacketInfo
	var packetInfoExist bool = false

	// check if we already get this message ( e.g. LOST of PUBREC  )
	packetInfoExist, packetInfo, _ = persistance.PacketExist("qos", persistance.PacketFindOpts{
		OriginClientID:  &ogiginalPacketID,
		OriginMessageID: &packet.MessageID,
	})

	//  Store package - For retry-purpose
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

	// we get our own messageID
	packet.MessageID = packetInfo.MessageID

	// Publish to all clients
	if err == nil && packetInfoExist == false {

		// publish packet to all subscribers
		var subscribed bool
		_, subscribed, err = broker.PublishPacket(packet, packet.Qos == 2)

		// nobody subscribe to this
		if subscribed == false {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")

			persistance.PacketDelete("qos", persistance.PacketFindOpts{
				OriginClientID:  &ogiginalPacketID,
				OriginMessageID: &packet.MessageID,
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
