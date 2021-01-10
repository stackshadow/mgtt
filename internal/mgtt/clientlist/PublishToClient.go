package clientlist

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// PublishToClient publish an message to an specific client ( and only to this ! )
//
// - if plugins are active, all plugins must accept this
//
// - the topic must match an subscription of the client
func PublishToClient(clientID string, topic string, payload []byte) (err error) {

	// find the client
	client := list[clientID]

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
