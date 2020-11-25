package acl

import (
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

// OnPublishRecvRequest write to broker
func checkACL(clientID string, username string, topic string, direction string) (allowed bool) {
	// if clientID is resend, this is an resended package... we allow this by default
	if clientID == "resend" {
		log.Debug().Str("topic", topic).Msg("This is an resendet packet, we allow it")
		return true
	}

	defer func() {
		if allowed == false {
			log.Warn().Str("topic", topic).Msg("Not allowed")
		}
	}()

	// if username is empty,
	if username == "" {
		username = "_anonym"
	}

	// try to get the acl
	entryArray, exist := config.Rules[username]
	if exist == false {
		return false
	}

	// iterate
	topicArray := strings.Split(topic, "/")
	for _, entry := range entryArray {

		if entry.Direction == direction {

			routeArray := strings.Split(entry.Route, "/")
			if client.Match(routeArray, topicArray) == true {
				allowed = entry.Allow
				return
			}
		}

	}

	return
}
