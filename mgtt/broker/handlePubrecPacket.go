package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

func (broker *Broker) handlePubrecPacket(event *Event) (err error) {

	// check package
	pubrecPacket, ok := event.packet.(*packets.PubrecPacket)
	if ok == false {
		err = errors.New("Package is not packets.PubrecPacket")
		return
	}

	// okay, we try to find the package in our "resend"-storage
	// find the package
	storedInfo, err := broker.retainedMessages.GetPacketByID("resend", pubrecPacket.MessageID)
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
		client.SendPubcomp(storedInfo.Packet.MessageID)
	} else {
		log.Info().Msg("No client stored, this should not happen, but we proceeed")
	}

	event.client.SendPubrel(pubrecPacket.MessageID)

	// err = broker.retainedMessages.DeletePacketWithID("unreleased", pubrecPacket.MessageID)

	return
}
