package broker

import "gitlab.com/mgtt/plugin"

func (broker *Broker) handleConackPacket(event *Event) (err error) {
	event.client.Connected = true
	broker.clients[event.client.ID()] = event.client

	// PLUGINS: call CallOnAcceptNewClient - check if we accept the client
	if err == nil {
		plugin.CallOnConnack(event.client.ID())
	}

	return
}
