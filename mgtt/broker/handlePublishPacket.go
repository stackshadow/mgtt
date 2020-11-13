package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handlePublishPacket(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	// call plugin
	acceptPublish := plugin.CallOnPublishRequest(client.ID(), client.Username(), packet.TopicName, string(packet.Payload))
	if acceptPublish == false {
		client.Close()
		return
	}

	// RETAINED-Packet
	if err == nil && packet.Retain == true && packet.Dup == false { // prevent multiple return

		// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
		// [MQTT-3.3.1-11] A zero byte retained message MUST NOT be stored as a retained message on the Server.
		if len(packet.Payload) == 0 {
			err = broker.retainedMessages.DeletePacketWithTopic("retained", packet.TopicName)
		} else {
			// [MQTT-3.3.1-5]
			err = broker.retainedMessages.StorePacketWithTopic("retained", packet.TopicName, packet)
		}

	}

	switch packet.Qos {
	case 0:
		err = broker.handlePublishPacketQoS0(client, packet)
	case 1:
		err = broker.handlePublishPacketQoS1(client, packet)
	case 2:
		err = broker.handlePublishPacketQoS2(client, packet)
	}

	return
}
