package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mgtt/persistance"
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

		// [MQTT-3.1.2-6] If CleanSession is set to 1, the Client and Server MUST discard any previous Session and start a new one.
		// This Session lasts as long as the Network Connection.
		// State data associated with this Session MUST NOT be reused in any subsequent Session.
		var sessionExist bool = false
		connectedClient.CleanSessionSet(packet.CleanSession)
		if packet.CleanSession == true {
			persistance.CleanSession(packet.ClientIdentifier)
		} else {
			sessionSubscriptions := persistance.SubscriptionsGet(packet.ClientIdentifier)
			connectedClient.SubScriptionsAdd(sessionSubscriptions)
			if len(sessionSubscriptions) > 0 {
				sessionExist = true
			}
		}

		// [MQTT-3.2.2-2]  If the Server accepts a connection with CleanSession set to 0, the value set in Session Present depends on
		// whether the Server already has stored Session state for the supplied client ID.
		// If the Server has stored Session state, it MUST set Session Present to 1 in the CONNACK packet.
		err = connectedClient.SendConnack(client.ConnackAccepted, sessionExist)

		// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
		if err == nil {
			plugin.CallOnConnected(connectedClient.ID())
		}
	}

	return
}
