package clientlist

/*
// ClientNew will add an client and create a new ID for it
func (br *broker.Broker) ClientNew(con net.Conn) (mgttClient *client.MgttClient, err error) {

	// create a new client with an new random-id
	guid := xid.New()
	mgttClient = client.New(con, broker.ConnectTimeout)
	mgttClient.IDSet(guid.String())

	err = b.ClientAdd(mgttClient)

	return
}

func (b *broker.Broker) ClientAdd(mgttClient *client.MgttClient) (err error) {
	if _, clientExist := broker.clients[mgttClient.ID()]; clientExist == true {
		err = errors.New("Client already exist")
	} else {
		broker.clients[mgttClient.ID()] = mgttClient
		log.Debug().Str("client", mgttClient.ID()).Msg("Client added")
	}
	return
}

// ClientIDSet will change an existing client
func (b *broker.Broker) ClientIDSet(clientID string, newClientID string) (err error) {

	// get the client
	client := broker.ClientGet(clientID)
	if client != nil {

		// remove the old client
		delete(broker.clients, clientID)

		// set
		client.IDSet(newClientID)

		// add the client
		err = broker.ClientAdd(client)
	} else {
		err = fmt.Errorf("Client with id '%s' not exist", newClientID)
	}

	return
}

// ClientGet return an mqtt-client if it exist with an clientID or nil if it not exist
func (b *broker.Broker) ClientGet(clientID string) (client *client.MgttClient) {

	if existingClient, exist := broker.clients[clientID]; exist == true {
		client = existingClient
	}

	return
}

func (b *broker.Broker) ClientRemove(clientID string) {
	delete(broker.clients, clientID)

	log.Debug().Str("client", clientID).Msg("Remove client from the list")
}
*/
