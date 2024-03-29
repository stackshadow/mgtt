package broker

import (
	"strings"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func (broker *Broker) onPacketPublish(client *client.MgttClient, packet *packets.PublishPacket) (err error) {

	// call plugin
	acceptPublish := plugin.CallOnPublishRequest(client.ID(), client.Username(), packet.TopicName)
	if acceptPublish == false {
		client.Close()
		return
	}

	// call plugin that possible handle the message
	if plugin.CallOnHandleMessage(client.ID(), packet.TopicName, packet.Payload) == true {
		return
	}

	// handle ping request
	if packet.TopicName == "$SYS/ping" {
		// construct the package
		pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
		pub.MessageID = 0
		pub.Retain = false
		pub.TopicName = "$SYS/pong"
		pub.Qos = 0
		pub.Payload = packet.Payload

		client.Publish(pub)
		return
	}

	// we not publish $SYS to others
	if strings.Index(packet.TopicName, "$SYS/") == 0 {
		return
	}

	// RETAINED-Packet
	// [MQTT-3.1.2.7] Retained messages do not form part of the Session state in the Server, they MUST NOT be deleted when the Session ends.
	if err == nil && packet.Retain == true && packet.Dup == false { // prevent multiple return

		// [MQTT-3.3.1-10] if payload is 0, an retained message MUST be removed
		// [MQTT-3.3.1-11] A zero byte retained message MUST NOT be stored as a retained message on the Server.
		if len(packet.Payload) == 0 {
			persistance.PacketDelete("retained",
				persistance.PacketFindOpts{
					Topic: &packet.TopicName,
				},
			)
		} else {

			persistance.PacketDelete("retained",
				persistance.PacketFindOpts{
					Topic: &packet.TopicName,
				},
			)

			// [MQTT-3.3.1-5]
			err = persistance.PacketStore("retained",
				&persistance.PacketInfo{
					Topic:   packet.TopicName,
					Payload: packet.Payload,
					Qos:     packet.Qos,
				},
			)

		}

	}

	switch packet.Qos {
	case 0:
		err = broker.onPacketPublishQoS0(client, packet)
	case 1:
		err = broker.onPacketPublishQoS1(client, packet)
	case 2:
		err = broker.onPacketPublishQoS1(client, packet)
	}

	return
}
