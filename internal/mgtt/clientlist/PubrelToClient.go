package clientlist

import (
	"errors"
)

// PubrelToClient send an PUBREL to an specific client
func PubrelToClient(clientID string, MessageID uint16) (err error) {

	// find the client
	listMutex.Lock()
	client := list[clientID]
	listMutex.Unlock()

	if client != nil {
		client.SendPubrel(MessageID)
	} else {
		err = errors.New("Client with id not exist")
	}

	return
}
