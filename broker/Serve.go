package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
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
		var serverURL *url.URL
		serverURL, err = url.Parse(config.URL)
		var serverListener net.Listener

		// check schema
		if serverURL.Scheme != "tcp" {
			err = fmt.Errorf("Unsupported scheme '%s'", serverURL.Scheme)
		}

		// start listener
		if err == nil {

			if config.CertFile == "" {
				serverListener, err = net.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port())

			} else {
				var cert tls.Certificate
				cert, err = tls.LoadX509KeyPair(config.CertFile, config.Keyfile)
				cfg := &tls.Config{
					Certificates:       []tls.Certificate{cert},
					InsecureSkipVerify: true,
				}
				serverListener, err = tls.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port(), cfg)
			}

		}

		if err == nil {
			if config.CertFile == "" {
				log.Info().Str("listen", serverURL.Host).Bool("tls", false).Msg("Listening")
			} else {
				log.Info().Str("listen", serverURL.Host).Bool("tls", true).Msg("Listening")
			}
		} else {
			log.Fatal().Err(err).Send()
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

					if err != nil {
						log.Error().Err(err).Send()
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
