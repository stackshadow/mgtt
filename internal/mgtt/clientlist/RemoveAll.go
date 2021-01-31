package clientlist

// RemoveAll will remove and disconnects all clients
func RemoveAll() {

	// prevent "fatal error: concurrent map iteration and map write"
	var keys []string

	// close all clients
	for clientID, client := range list {
		client.Close()
		keys = append(keys, clientID)
	}

	// Delete all clints
	for _, key := range keys {
		delete(list, key)
	}

}
