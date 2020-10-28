package broker

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handleDisConnectPacket(connectedClient *client.MgttClient) (err error) {

	log.Info().Str("clientid", connectedClient.ID()).Msg("Remove client from client-list")
	delete(broker.clients, connectedClient.ID())

	return
}
