package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

func (broker *Broker) handleConackPacket(connectedClient *client.MgttClient, packet *packets.ConnackPacket) (err error) {
	connectedClient.Connected = true
	broker.clients[connectedClient.ID()] = connectedClient

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	if err == nil {
		plugin.CallOnConnack(connectedClient.ID())
	}

	return
}
