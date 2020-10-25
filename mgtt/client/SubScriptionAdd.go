package client

// SubScriptionAdd will add a single subscription-topic to the client
func (client *MgttClient) SubScriptionAdd(topic string) {

	// build subscription map to dedublicate subscriptions
	subsscriptionmap := make(map[string]bool)
	for _, subscription := range client.subscriptionTopics {
		subsscriptionmap[subscription] = true
	}
	subsscriptionmap[topic] = true

	// recreate array and set it back
	var newSubscriptions []string
	for subsscription := range subsscriptionmap {
		newSubscriptions = append(newSubscriptions, subsscription)
	}

	client.subscriptionTopics = newSubscriptions
}

// SubScriptionsAdd will add multiple subscription-topics to the client
func (client *MgttClient) SubScriptionsAdd(topics []string) {

	// build subscription map to dedublicate subscriptions
	subsscriptionmap := make(map[string]bool)
	for _, subscription := range client.subscriptionTopics {
		subsscriptionmap[subscription] = true
	}
	for _, subscription := range topics {
		subsscriptionmap[subscription] = true
	}

	// recreate array and set it back
	var newSubscriptions []string
	for subsscription := range subsscriptionmap {
		newSubscriptions = append(newSubscriptions, subsscription)
	}

	client.subscriptionTopics = newSubscriptions

}
