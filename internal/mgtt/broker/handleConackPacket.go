package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func (broker *Broker) handleConackPacket(connectedClient *client.MgttClient, packet *packets.ConnackPacket) (err error) {
	connectedClient.Connected = true

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	if err == nil {
		plugin.CallOnConnack(connectedClient.ID())
	}

	return
}
