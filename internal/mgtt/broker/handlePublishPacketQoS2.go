package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	messagestore "gitlab.com/mgtt/internal/mgtt/messageStore"
)

func (broker *Broker) handlePublishPacketQoS2(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	var options messagestore.StoreResendPacketOptions
	options.OriginID = packet.MessageID

	//  QoS2 - Store package only if its not duplicated
	if packet.Dup == false {
		// we need a new ID
		broker.lastIDLock.Lock()
		options.ClientID = client.ID()
		options.ResendAt = time.Now().Add(time.Minute * 1)
		options.Packet = packet

		// store the packet to unreleased
		err = broker.retainedMessages.StoreResendPacket("unreleased", &broker.lastID, &options)

		broker.lastIDLock.Unlock()
	}

	// send pubrec
	if err == nil && packet.Dup == false {
		client.SendPubrec(options.OriginID)
	}

	return
}
