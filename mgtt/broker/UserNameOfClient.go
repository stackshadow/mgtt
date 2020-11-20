package broker

 

// UserNameOfClient return the username of an client
func (broker *Broker) UserNameOfClient(clientID string ) ( username string ) {

	// find the client
	for _, curClient := range broker.clients {
		if curClient.ID() == clientID {
			username = curClient.Username()
			return
		}
	}

	return
}
