package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func (broker *Broker) onPacketConnect(connectedClient *client.MgttClient, packet *packets.ConnectPacket) (err error) {

	// MQTT-3.1.0-2
	// Check if the client is already connected
	if err == nil { // prevent multiple return
		if exists := clientlist.Exist(packet.ClientIdentifier); exists == true {
			err = errors.New("Protocol violation. Client already exist")
		}
	}

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	if err == nil { // prevent multiple return
		accepted := plugin.CallOnAcceptNewClient(connectedClient.ID(), packet.Username, string(packet.Password))
		if accepted == false {
			err = connectedClient.SendConnack(client.ConnackUnauthorized)
			err = errors.New("Client not accepted by plugin")
		}
	}

	// add client to the list
	if err == nil { // prevent multiple return

		// Move the client to the newID
		clientlist.Move(connectedClient.ID(), packet.ClientIdentifier)

		// store the username
		connectedClient.UsernameSet(packet.Username)

		// set the client to connected so that the broker will accept other packets from it
		connectedClient.Connected = true

		// reset timeout
		connectedClient.ResetTimeout()

		// Las will message ?
		if packet.WillFlag == true {
			// we create a new publish packet
			pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
			pubPacket.Retain = packet.WillRetain
			pubPacket.TopicName = packet.WillTopic
			pubPacket.Payload = packet.WillMessage
			pubPacket.Qos = packet.WillQos

			connectedClient.LastWillSet(pubPacket)
		}

		// send CONACK
		err = connectedClient.SendConnack(client.ConnackAccepted)

		// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
		if err == nil {
			plugin.CallOnConnected(connectedClient.ID())
		}
	}

	return
}
