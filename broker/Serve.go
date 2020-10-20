package broker

import (
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// Serve will create a new broker and wait for clients
func Serve(config Config) (broker *Broker, err error) {

	broker = &Broker{
		clients:      make(map[string]*client.MgttClient),
		clientEvents: make(chan *Event, 10),
		pubrec:       make(map[uint16]*messagestore.StoreResendPacketOption),
	}

	// retainedMessages-db
	broker.retainedMessages, err = messagestore.Open("messages.db")
	if err != nil {
		return
	}

	// incoming clients
	go func() {
		var serverListener net.Listener
		serverListener, err = net.Listen("tcp", config.URL)
		if err != nil {
			log.Error().Err(err).Send()
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

					for {

						// new event
						newEvent := Event{
							client: newClient,
						}

						// wait for a packet
						newEvent.packet, err = newClient.ReadPacket()
						if err != nil {
							break
						}

						broker.clientEvents <- &newEvent
					}

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
			log.Debug().Msg("Check if packets need to be resended")
			broker.retainedMessages.IterateResendPackets("resend", func(storedInfo *messagestore.StoreResendPacketOption) {

				// check if time is up
				if time.Now().After(storedInfo.ResendAt) == true {

					storedInfo.Packet.Dup = true                             // this is an duplicate packet
					storedInfo.Packet.Retain = false                         // resend, not retain ;)
					storedInfo.Packet.MessageID = storedInfo.BrokerMessageID // we use our message ID

					log.Debug().
						Uint16("packet-mid", storedInfo.Packet.MessageID).
						Uint16("broker-mid", storedInfo.BrokerMessageID).
						Str("topic", storedInfo.Packet.TopicName).
						Msg("Resend packet")

					// new event
					newEvent := Event{
						client: nil,
						packet: storedInfo.Packet,
					}

					broker.clientEvents <- &newEvent
				}

			})

		}
	}()

	return
}
