package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
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
			log.Info().Str("listen", serverURL.Host).
				Bool("tls", false).
				Msg("Connected")
		} else {
			log.Info().Str("listen", serverURL.Host).
				Bool("tls", true).
				Str("cert", config.CertFile).
				Str("key", config.KeyFile).
				Msg("Listening")
		}
	} else {
		log.Fatal().Err(err).Send()
	}

	// create the client
	newClient := client.New(clientListener, ConnectTimeout)
	newClient.SubScriptionAdd("#")

	// run communicate
	newClient.Communicate()

	// We need a random uuid
	var newUUID uuid.UUID
	newUUID, err = uuid.NewRandom()
	if err != nil {
		return
	}

	// send it to the client
	err = newClient.SendConnect(username, password, newUUID.String())
	if err != nil {
		return
	}

	// do communication
	var normalClose bool
	for {

		// get packet from the client-buffer
		recvdPacket := newClient.GetPacket()

		// if we get a nil-packet, client-connection is closed
		if recvdPacket == nil {
			err = nil
			break
		}

		normalClose, err = broker.loopHandleClientPackets(newClient, recvdPacket)
		if err != nil || normalClose {
			break
		}
	}

	return
}
