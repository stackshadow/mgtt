package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/plugin"
)

// PublishPacket publish a packet to all subscribers
func (broker *Broker) PublishPacket(packet *packets.PublishPacket, once bool) (messagedelivered bool, err error) {

	var published bool

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
	// because it matches an established subscription
	packet.Retain = false

	// PLUGINS: call CallOnPublishRequest - check if publish is accepted
	for _, client := range broker.clients {
		clientID := client.ID()
		userName := client.Username()
		if plugin.CallOnSubscribeRequest(clientID, userName, packet.TopicName) == true {
			published, err = client.Publish(packet)
			if once == true {
				if published == true {
					return true, nil
				}
			}
			messagedelivered = messagedelivered || published
		}
	}

	return
}
