package acl

import "github.com/rs/zerolog/log"

// OnSendToSubscriberRequest get called when the broker try to publish a message to an subscriber
//
// if an plugin set it to false, the message will NOT be sended to clientID
//
// This function gets called BEFORE check if the subscriber subscribe to the topic
//
// clientID: The clientID of the target client
// username: The username of the target client
// publishTopic: The topic the broker try to publish to the subscriber
func OnSendToSubscriberRequest(clientID string, username string, publishTopic string) (accepted bool) {

	// check global permission
	allowedGlobally := checkACL(clientID, "_global", publishTopic, "r")

	// check for specific username
	accepted = checkACL(clientID, username, publishTopic, "r") || allowedGlobally

	if accepted == false {
		log.Warn().Str("topic", publishTopic).Msg("Not allowed")
	}

	return accepted
}
