package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

// Serve will create a new broker and wait for clients
//
// - this create also an retained message with topic "$SYS/broker/version" that contains the version
func (b *Broker) Serve(config Config) (err error) {

	err = persistance.Open(config.DBFilename)

	// Delete Broker-version if exist
	brokerVersionTopic := "$SYS/broker/version"
	persistance.PacketDelete("retained",
		persistance.PacketFindOpts{
			Topic: &brokerVersionTopic,
		},
	)
	// Set the broker-version
	persistance.PacketStore("retained",
		&persistance.PacketInfo{
			Topic:   brokerVersionTopic,
			Payload: []byte(config.Version),
		},
	)

	var serverURL *url.URL
	serverURL, err = url.Parse(config.URL)

	// check schema
	if serverURL.Scheme != "tcp" {
		err = fmt.Errorf("Unsupported scheme '%s'", serverURL.Scheme)
	}

	// check if an error occurred
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// non-tls
	if config.TLS == true {
		var TLSConfig *tls.Config
		TLSConfig, err = getTLSConfig(config)
		if err == nil {
			b.serverListener, err = tls.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port(), TLSConfig)
		}
	} else {
		b.serverListener, err = net.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port())
	}

	// check if an error occured
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// logging
	if config.CertFile == "" {
		log.Info().Str("listen", serverURL.Host).
			Bool("tls", false).
			Msg("Listening")
	} else {
		log.Info().Str("listen", serverURL.Host).
			Bool("tls", true).
			Str("ca", config.CAFile).
			Str("cert", config.CertFile).
			Str("key", config.KeyFile).
			Msg("Listening")
	}

	// retry TODO
	// go b.loopHandleResendPackets()

	for {

		// wait for a new client
		log.Info().Msg("Wait for new client")
		var newConnection net.Conn
		newConnection, err = b.serverListener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Accept()")
			break
		}

		// handle a new client
		if err == nil {
			go handleNewClient(newConnection)
		}
	}

	return
}
