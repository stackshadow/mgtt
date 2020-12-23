package clientlist

import (
	"errors"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func Add(existingClient *client.MgttClient) (err error) {

	if _, clientExist := list[existingClient.ID()]; clientExist == true {
		err = errors.New("Client already exist")
	} else {
		list[existingClient.ID()] = existingClient
		log.Debug().Str("client", existingClient.ID()).Msg("Client added")
	}

	return
}
