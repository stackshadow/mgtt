package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

// PublishToClient publish an message to an specific client ( and only to this ! )
func (broker *Broker) PublishToClient(clientID string, topic string, payload []byte) (err error) {

	// find the client
	var client clientlist.Client = clientlist.Get(clientID)

	if client != nil {
		// construct the package
		pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
		pub.MessageID = 0
		pub.Retain = false
		pub.TopicName = topic
		pub.Payload = payload
		pub.Qos = 0

		client.Publish(pub)
	} else {
		err = errors.New("Client with id not exist")
	}

	return
}
