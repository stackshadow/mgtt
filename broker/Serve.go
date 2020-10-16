package broker

import (
	"net"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
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

	// incoming clients
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

	// retrys
	go func() {

		for {

			// wait a bit
			time.Sleep(time.Minute * 1)

			// check if we need to resend messages that are not replyed with PUBACK
			log.Debug().Msg("Check if we packets we should resend")
			broker.retainedMessages.IteratePackets("resend", func(retainedPacket *packets.PublishPacket) {
				for _, client := range broker.clients {

					// [MQTT-3.3.1-8]
					// When sending a PUBLISH Packet to a Client the Server MUST set the RETAIN flag to 1
					// if a message is sent as a result of a new subscription being made by a Client.
					retainedPacket.Retain = true

					client.Publish(retainedPacket)
				}
			})

		}
	}()

	return
}
