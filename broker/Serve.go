package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

// Serve will create a new broker and wait for clients
func (broker *Broker) Serve(config Config) (err error) {

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

			// non-tls
			if config.CertFile == "" {
				serverListener, err = net.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port())

			} else { // tls
				var cert tls.Certificate
				cert, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
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
				log.Info().Msg("New client connected")

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
	go broker.loopHandleResendPackets()

	// handle packets
	broker.loopHandleBrokerPackets()

	return
}
