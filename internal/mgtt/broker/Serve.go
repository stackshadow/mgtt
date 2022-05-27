package broker

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/mgtt/internal/mgtt/persistance"
	"gitlab.com/mgtt/internal/mgtt/server"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

// Serve will create a new broker and wait for clients in a goroutine
//
// - this create also an retained message with topic "$METRIC/broker/version" that contains the version
func (b *Broker) Serve() (done chan bool, err error) {

	persistance.MustOpen(config.Globals.DB)

	// Delete Broker-version if exist
	brokerVersionTopic := "$METRIC/broker/version"
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
	serverURL, err = url.Parse(config.Globals.URL)
	utils.PanicOnErr(err)

	// check schema
	if serverURL.Scheme != "tcp" {
		err = fmt.Errorf("Unsupported scheme '%s'", serverURL.Scheme)
		utils.PanicOnErr(err)
	}

	// create a server
	serverListener := server.Create(serverURL.Hostname(), serverURL.Port())
	serverListener.MustInit(config.Globals.TLS.CA.File, config.Globals.TLS.Cert.File)

	// retry
	go b.loopHandleResendPackets()

	go func() {
		for {

			// wait for a new client
			log.Info().Msg("Wait for new client")
			newClient := serverListener.Accept()
			err = clientlist.Add(newClient)

			// handle a new client
			if err == nil {
				go handleNewClient(b, newClient)
			}
		}

		done <- true
	}()

	return
}
