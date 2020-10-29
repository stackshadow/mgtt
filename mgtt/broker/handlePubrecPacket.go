package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePubrecPacket(connectedClient *client.MgttClient, packet *packets.PubrecPacket) (err error) {

	// okay, we try to find the package in our "resend"-storage
	// find the package
	storedInfo, err := broker.retainedMessages.GetPacketByID("resend", packet.MessageID)
	if err != nil {
		return err
	}

	// cool, we look for the client
	if storedInfo.ClientID != "" {

		client, clientExist := broker.clients[storedInfo.ClientID]
		if clientExist == false {
			err = errors.New("Client '" + storedInfo.ClientID + "' not exist anymore.")
		}

		// send pubcomp
		client.SendPubcomp(storedInfo.OriginID)
	} else {
		log.Info().Msg("No client stored, this should not happen, but we proceeed")
	}

	connectedClient.SendPubrel(packet.MessageID)

	// err = broker.retainedMessages.DeletePacketWithID("unreleased", pubrecPacket.MessageID)

	return
}
