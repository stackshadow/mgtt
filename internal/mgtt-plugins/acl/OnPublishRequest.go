package acl

import (
	"fmt"

	"gitlab.com/mgtt/internal/mgtt/broker"
)

// OnPublishRequest get called when an publisher try to publish to the broker
func OnPublishRequest(clientID string, username string, topic string) (accepted bool) {

	allowed := checkACL(clientID, username, topic, "w")

	if allowed == false {
		broker.Current.PublishToClient(
			clientID,
			"$SYS/self/error",
			[]byte(fmt.Sprintf("Access to '%s' denied", topic)),
		)
	}

	return allowed
}
