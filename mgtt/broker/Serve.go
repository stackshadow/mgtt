package broker

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/cli"
	"gitlab.com/mgtt/client"
)

// Serve will create a new broker and wait for clients
func (broker *Broker) Serve(config Config) (err error) {

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
			log.Info().Str("listen", serverURL.Host).
				Bool("tls", true).
				Str("cert", config.CertFile).
				Str("key", config.KeyFile).
				Msg("Listening")
		}
	} else {
		log.Fatal().Err(err).Send()
	}

	// retrys
	var retryPackets chan packets.ControlPacket = make(chan packets.ControlPacket, 100)
	broker.loopHandleResendPackets(retryPackets)
	go func() {

		netserver, _ := net.Pipe()
		retryClient := client.New(netserver, 0)
		retryClient.IDSet("resend")
		retryClient.Connected = true

		for {
			retryPacket := <-retryPackets
			normalClose, err := broker.loopHandleBrokerPacket(retryClient, retryPacket)
			if err != nil || normalClose {
				break
			}

			// a small delay to not flood our clients
			time.Sleep(time.Millisecond * 500)
		}

		netserver.Close()
	}()

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
				newClient := client.New(newConnection, cli.CLI.ConnectTimeout)
				log.Info().Msg("New client connected")

				// run communicate
				newClient.Communicate()

				// do communication
				var normalClose bool
				for {
					recvdPacket := newClient.GetPacket()
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
