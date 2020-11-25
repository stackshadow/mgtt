package broker

import "github.com/eclipse/paho.mqtt.golang/packets"

// Publish will publish a message to all clients
func (broker *Broker) Publish(topic string, payload []byte, retain bool, QoS byte) (err error) {

	// construct the package
	pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.MessageID = broker.lastID
	pub.Retain = retain
	pub.TopicName = topic
	pub.Payload = payload
	pub.Qos = QoS
	broker.lastID++

	for _, client := range broker.clients {
		client.Publish(pub)
	}

	return
}
