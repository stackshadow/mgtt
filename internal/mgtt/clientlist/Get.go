package clientlist

import "gitlab.com/mgtt/internal/mgtt/client"

// Get will return an mgttClient from an given clientID
func Get(clientID string) (existingClient *client.MgttClient) {
	existingClient, _ = list[clientID]
	return
}
