package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleConnectPacket(connectedClient *client.MgttClient, packet *packets.ConnectPacket) (err error) {

	// set the client id
	connectedClient.IDSet(packet.ClientIdentifier)

	// MQTT-3.1.0-2
	// Check if the client is already connected
	if err == nil { // prevent multiple return
		if _, exists := broker.clients[connectedClient.ID()]; exists == true {
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

		// store the username
		connectedClient.UsernameSet(packet.Username)

		// set the client to connected so that the broker will accept other packets from it
		connectedClient.Connected = true

		// reset timeout
		connectedClient.ResetTimeout()

		// add client to the list
		log.Info().Str("clientid", connectedClient.ID()).Msg("Add new client to client-list")
		broker.clients[connectedClient.ID()] = connectedClient

		// send CONACK
		err = connectedClient.SendConnack(client.ConnackAccepted)

		// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
		if err == nil {
			plugin.CallOnConnected(connectedClient.ID())
		}
	}

	return
}
