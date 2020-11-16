package acl

// OnPublishRequest get called when an publisher try to publish to the broker
func OnPublishRequest(clientID string, username string, topic string) (accepted bool) {
	return checkACL(clientID, username, topic, "w")
}
