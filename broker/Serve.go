package broker

import (
	"net"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// Serve will create a new broker and wait for clients
func Serve(config Config) (broker *Broker, err error) {

	broker = &Broker{
		clients:       make(map[string]*client.MgttClient),
		clientEvents:  make(chan *client.Event, 10),
		pendingEvents: make(map[uint16]*client.Event),
	}

	// retainedMessages-db
	broker.retainedMessages, err = messagestore.Open()
	if err != nil {
		return
	}

	go func() {
		var serverListener net.Listener
		serverListener, err = net.Listen("tcp", config.URL)
		if err != nil {
			return
		}
		for {

			// wait for a new client
			log.Info().Msg("Wait for new client")
			var newConnection net.Conn
			newConnection, err = serverListener.Accept()
			if err != nil {
				log.Error().Err(err).Msg("not accepted")
			}

			// create a new client
			if err == nil {
				newClient := client.New(newConnection)
				log.Info().Str("clientid", newClient.ID()).Msg("New client connected")

				// do communication
				go func() {
					newClient.Communicate(broker.clientEvents)

					log.Info().Str("clientid", newClient.ID()).Msg("Remove client from client-list")
					delete(broker.clients, newClient.ID())
				}()
			}
		}
	}()

	return
}
