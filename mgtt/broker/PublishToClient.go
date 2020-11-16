package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

// PublishToClient publish an message to an specific client ( and only to this ! )
func (broker *Broker) PublishToClient(clientID string, topic string, payload []byte) (err error) {

	// find the client
	var client *client.MgttClient
	for _, curClient := range broker.clients {
		if curClient.ID() == clientID {
			client = curClient
		}
	}

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
