package broker

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) loopReadPackets(Client *client.MgttClient) (err error) {

	for {

		// new event
		newEvent := Event{
			client: Client,
		}

		// wait for a packet
		newEvent.packet, err = Client.ReadPacket()
		if err != nil {
			break
		}

		broker.clientEvents <- &newEvent
	}

	log.Info().
		Str("clientid", Client.ID()).
		Err(err).
		Msg("Remove client from client-list")

	Client.Close()
	delete(broker.clients, Client.ID())

	return
}
