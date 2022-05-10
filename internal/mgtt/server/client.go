package server

import (
	"net"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

// Accept will block until a newClient connected to the server
func (l *Listener) Accept() (newClient *client.MgttClient) {

	var err error

	log.Info().Msg("Wait for new client")
	var newConnection net.Conn
	newConnection, err = l.listener.Accept()
	utils.PanicOnErr(err)

	// create a new client
	newClient = &client.MgttClient{}

	// init it
	newClient.Init(newConnection, config.Values.Timeout)

	return
}
