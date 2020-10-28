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
)

// Connect will connect to an broker
func (broker *Broker) Connect(config Config, username, password string) (err error) {

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
	newClient := client.New(clientListener, cli.CLI.ConnectTimeout)
	newClient.SubScriptionAdd("#")

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

	broker.loopReadPackets(newClient)

	return
}
