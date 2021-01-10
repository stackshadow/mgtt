package clientlist

// Get will return an mgttClient from an given clientID
func Get(clientID string) (existingClient Client) {
	existingClient, _ = list[clientID]
	return
}
