package clientlist

// RemoveAll will remove and disconnects all clients
func RemoveAll() {

	for clientID, client := range list {
		client.Close()
		delete(list, clientID)
	}

}
