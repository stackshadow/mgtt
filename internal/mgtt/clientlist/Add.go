package clientlist

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func Add(existingClient Client) (err error) {

	// mutex
	listMutex.Lock()
	defer listMutex.Unlock()

	if _, clientExist := list[existingClient.ID()]; clientExist == true {
		err = errors.New("Client already exist")
	} else {
		list[existingClient.ID()] = existingClient
		log.Debug().Str("client", existingClient.ID()).Msg("Client added")
	}

	return
}
