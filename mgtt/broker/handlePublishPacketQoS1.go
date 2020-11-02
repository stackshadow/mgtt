package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) handlePublishPacketQoS1(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var packetID uint16 = packet.MessageID

	//  QoS1 - Store package only if its not duplicated
	if packet.Dup == false {

		// we need a new ID
		broker.lastIDLock.Lock()

		options := messagestore.StoreResendPacketOptions{
			ClientID: client.ID(),
			ResendAt: time.Now().Add(time.Minute * 1),
			Packet:   packet,
		}

		//
		err = broker.retainedMessages.StoreResendPacket("resend", &broker.lastID, &options)

		broker.lastIDLock.Unlock()
	}

	// Publish to all clients
	var messagedelivered bool
	if err == nil {

		// publish packet to all subscribers
		messagedelivered, err = broker.PublishPacket(packet, false)

		// no message delivered
		if messagedelivered == true {
			broker.retainedMessages.DeletePacketWithID("resend", packet.MessageID)
		} else {
			log.Info().
				Str("topic", packet.TopicName).
				Uint16("mid", packet.MessageID).
				Msg("Nobody is interested in this message")

		}

	}

	client.SendPuback(packetID)
	return
}
