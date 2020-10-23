package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/cli"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

// Connect will connect to an broker
func Connect(config Config, username, password string) (broker *Broker, err error) {

	broker = &Broker{
		clients:      make(map[string]*client.MgttClient),
		clientEvents: make(chan *Event, 10),
	}

	// retainedMessages-db
	broker.retainedMessages, err = messagestore.Open(cli.CLI.DBFilename)
	if err != nil {
		return
	}

	// connect to the broker
	var serverURL *url.URL
	serverURL, err = url.Parse(config.URL)
	var clientListener net.Conn

	// check schema
	if serverURL.Scheme != "tcp" {
		err = fmt.Errorf("Unsupported scheme '%s'", serverURL.Scheme)
	}

	// connect
	if err == nil {

		// non-tls
		if config.CertFile == "" {
			clientListener, err = net.Dial("tcp", serverURL.Hostname()+":"+serverURL.Port())

		} else { // tls
			var cert tls.Certificate
			cert, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
			cfg := &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			}
			clientListener, err = tls.Dial("tcp", serverURL.Hostname()+":"+serverURL.Port(), cfg)
		}

	}

	if err == nil {
		if config.CertFile == "" {
			log.Info().Str("listen", serverURL.Host).Bool("tls", false).Msg("Connected")
		} else {
			log.Info().Str("listen", serverURL.Host).Bool("tls", true).Msg("Connected")
		}
	} else {
		log.Fatal().Err(err).Send()
	}

	// create the client
	newClient := client.New(clientListener)

	// retrys
	go broker.loopHandleResendPackets()

	// handle client packets
	go broker.loopHandleClientPackets()

	// Send Connect-packet
	var newUUID uuid.UUID
	newUUID, err = uuid.NewRandom()
	if err != nil {
		return
	}

	// add the client to

	err = newClient.SendConnect(username, password, newUUID.String())
	if err != nil {
		return
	}

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

	log.Info().Str("clientid", newClient.ID()).Msg("Disconnected")

	return
}
