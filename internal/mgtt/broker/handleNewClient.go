package broker

import (
	"net"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func handleNewClient(newConnection net.Conn) {

	var err error
	var recvdPacket packets.ControlPacket

	newClient := client.New(newConnection, ConnectTimeout)
	err = clientlist.Add(newClient)

	if err == nil {

		// inform the plugins
		plugin.CallOnNewClient(newClient.RemoteAddr())

		// do communication
		var normalClose bool
		for {

			// get packet from the client-buffer
			recvdPacket, err = newClient.PacketRead()

			// if we get an error
			if err != nil {
				break
			}

			// handle the packet was broker
			normalClose, err = Current.handlePacketsForBroker(newClient, recvdPacket)
			if err != nil || normalClose == true {
				break
			}

		}

		// log-error
		if err != nil {
			log.Error().Err(err).Send()
		}

		// last-Will-message
		if lastWillPacket := newClient.LastWillGet(); lastWillPacket != nil {
			Current.handlePublishPacket(newClient, lastWillPacket)
		}

		// Remove the client from the list
		clientlist.Remove(newClient.ID())

	}
}
