package clientlist

// Get will return an mgttClient from an given clientID
func Get(clientID string) (existingClient Client) {
	// mutex
	listMutex.Lock()
	defer listMutex.Unlock()

	existingClient, _ = list[clientID]
	return
}
