package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleConnectPacket(event *client.Event) (err error) {

	// check package
	packet, ok := event.Packet.(*packets.ConnectPacket)
	if ok == false {
		err = errors.New("Package is not packets.ConnectPacket")
	}

	// set the client id
	event.Client.IDSet(packet.ClientIdentifier)

	// MQTT-3.1.0-2
	// Check if the client is already connected
	if err == nil { // prevent multiple return
		if _, exists := broker.clients[event.Client.ID()]; exists == true {
			err = errors.New("Protocol violation")
		}
	}

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	if err == nil { // prevent multiple return
		accepted := plugin.CallOnAcceptNewClient(event.Client.ID(), packet.Username, string(packet.Password))
		if accepted == false {
			err = event.Client.SendConnack(client.ConnackUnauthorized)
			err = errors.New("Client not accepted by plugin")
		}
	}

	// add client to the list
	if err == nil { // prevent multiple return
		log.Info().Str("clientid", event.Client.ID()).Msg("Add new client to client-list")
		broker.clients[event.Client.ID()] = event.Client

		// send CONACK
		err = event.Client.SendConnack(client.ConnackAccepted)
	}

	return
}
