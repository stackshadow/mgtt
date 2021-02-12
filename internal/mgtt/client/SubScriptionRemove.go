package client

// SubScriptionRemove will remove an topic from the subscriptionlist
func (client *MgttClient) SubScriptionRemove(topic string) {

	// the new list
	var newSubscriptionList []string

	// add subscriptions to the new list
	for _, subscription := range client.subscriptionTopics {
		if subscription != topic {
			newSubscriptionList = append(newSubscriptionList, subscription)
		}
	}

	// store it
	client.subscriptionTopics = newSubscriptionList
}

// SubScriptionsRemove will remove all topics
func (client *MgttClient) SubScriptionsRemove(topics []string) {

	for _, subscription := range topics {
		client.SubScriptionRemove(subscription)
	}

}
