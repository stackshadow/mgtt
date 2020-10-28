package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
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
		newEvent.packet = Client.GetPacket()
		if err != nil {
			break
		}

		// CONNACK-Packet
		switch recvPacket := newEvent.packet.(type) {
		case *packets.ConnectPacket:
			err = broker.handleConnectPacket(Client, recvPacket)
			if err != nil {
				break
			}
			continue
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
