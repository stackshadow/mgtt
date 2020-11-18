package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
	"gitlab.com/mgtt/plugin"
)

// Serve will create a new broker and wait for clients
func (broker *Broker) Serve(config Config) (err error) {

	// retainedMessages-db
	broker.retainedMessages, err = messagestore.Open(config.DBFilename)
	if err != nil {
		return
	}

	//	store version information
	pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.MessageID = 0
	pub.Retain = false
	pub.TopicName = "$SYS/broker/version"
	pub.Payload = []byte(config.Version)
	pub.Qos = 0
	broker.retainedMessages.StorePacketWithTopic("retained", pub.TopicName, pub)

	var serverURL *url.URL
	serverURL, err = url.Parse(config.URL)
	var serverListener net.Listener

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
			serverListener, err = tls.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port(), TLSConfig)
		}
	} else {
		serverListener, err = net.Listen("tcp", serverURL.Hostname()+":"+serverURL.Port())
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

	// retry
	go broker.loopHandleResendPackets()

	for {

		// wait for a new client
		log.Info().Msg("Wait for new client")
		var newConnection net.Conn
		newConnection, err = serverListener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Accept()")
		}

		// create a new client
		if err == nil {

			go func() {
				newClient := client.New(newConnection, ConnectTimeout)
				log.Info().Msg("New client connected")

				plugin.CallRemoteAddr(newClient.RemoteAddr())

				// run communicate
				newClient.Communicate()

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

					normalClose, err = broker.loopHandleBrokerPacket(newClient, recvdPacket)
					if err != nil || normalClose {
						break
					}
				}

				if err != nil {
					log.Error().Err(err).Send()
				}
				broker.handleDisConnectPacket(newClient)
				newClient.Close()

			}()

		}
	}

}
