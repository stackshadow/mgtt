package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleConnectPacket(event *client.Event) (err error) {
	packet, ok := event.Packet.(*packets.ConnectPacket)
	if ok == false {
		log.Error().Str("clientid", event.Client.ID()).Msg("Expected ConnectPacket")
		return
	}

	// set the client id
	event.Client.IDSet(packet.ClientIdentifier)

	log.Debug().
		Str("clientid", event.Client.ID()).
		Msg("RCV ConnectPacket")

	// MQTT-3.1.0-2
	// Check if the client is already connected
	if _, exists := broker.clients[event.Client.ID()]; exists == true {
		err = errors.New("Protocol violation")
		return
	}

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	accepted := plugin.CallOnAcceptNewClient(event.Client.ID(), packet.Username, string(packet.Password))
	if accepted == false {
		log.Error().Str("clientid", event.Client.ID()).Msg("Client not accepted by plugin")
		event.Client.SendConnack(client.ConnackUnauthorized)
		return
	}

	// add client to the list
	log.Info().Str("clientid", event.Client.ID()).Msg("Add new client to client-list")
	broker.clients[event.Client.ID()] = event.Client

	// send CONACK
	event.Client.SendConnack(client.ConnackAccepted)

	return
}
