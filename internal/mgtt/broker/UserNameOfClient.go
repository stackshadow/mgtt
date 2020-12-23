package broker

import "gitlab.com/mgtt/internal/mgtt/clientlist"

// UserNameOfClient return the username of an client
func (broker *Broker) UserNameOfClient(clientID string) (username string) {

	// find the client
	client := clientlist.Get(clientID)
	if client != nil {
		username = client.Username()
	}

	return
}
