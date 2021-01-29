package client

// Subscriptions return the current subscriptions of an client
func (client *MgttClient) Subscriptions() []string {
	return client.subscriptionTopics
}
