package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePublishPacket(event *client.Event) (err error) {

	// check package
	packet, ok := event.Packet.(*packets.PublishPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	if err == nil { // prevent multiple return
		// retain message ?
		if packet.Retain == true {

			// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
			// [MQTT-3.3.1-11] A zero byte retained message MUST NOT be stored as a retained message on the Server.
			if len(packet.Payload) == 0 {
				err = broker.retainedMessages.DeletePacketWithTopic("retained", packet.TopicName)
			} else {

				// [MQTT-3.3.1-5]
				err = broker.retainedMessages.StorePacket("retained", packet)
			}
		}
	}

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a
	// Client because it matches an established subscription
	if err == nil { // prevent multiple return
		packet.Retain = false
		for _, client := range broker.clients {
			err = client.Publish(packet)
		}
	}

	// Handle QoS-1 - Acknowledged delivery
	if packet.Qos == client.SubackQoS1 {
		// we ignore the returned err by purpose
		event.SendPuback()
	}


	return
}
