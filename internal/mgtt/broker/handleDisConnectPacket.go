package broker

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handleDisConnectPacket(connectedClient *client.MgttClient) (err error) {
	log.Info().Str("client", connectedClient.ID()).Msg("Disconnect received")
	return
}
