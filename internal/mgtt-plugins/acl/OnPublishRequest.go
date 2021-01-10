package acl

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

// OnPublishRequest get called when an publisher try to publish to the broker
func OnPublishRequest(clientID string, username string, topic string) (accepted bool) {
	var err error

	allowed := checkACL(clientID, username, topic, "w")

	if allowed == false {
		err = clientlist.PublishToClient(
			clientID,
			"$SYS/self/error",
			[]byte(fmt.Sprintf("Access to '%s' denied", topic)),
		)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}

	return allowed
}
