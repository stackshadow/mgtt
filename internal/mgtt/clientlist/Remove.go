package clientlist

import "github.com/rs/zerolog/log"

// Remove will remove an client with the given clientID and disconnects it
func Remove(clientID string) {
	if client, exist := list[clientID]; exist == true {
		client.Close()
		delete(list, clientID)
	} else {
		log.Warn().Str("client", clientID).Msg("Can not remove client from the list, it not exist on the list")
	}

}
