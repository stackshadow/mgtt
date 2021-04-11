package auth

import (
	"strings"

	"gitlab.com/mgtt/internal/mgtt/client"
)

var adminTopicsEnabled bool

// OnHandleMessage gets called after OnPublishRequest
//
// If this function return true, the plugin handled the message and no other plugin will get it
//
// If a plugin handle the message, it will NOT sended to subscribers
func OnHandleMessage(originClientID string, topic string, payload []byte) (handled bool) {

	// who is currently logged in
	if topic == "$SYS/self/user/get" {
		handled = true
		go onSelfUserGet(originClientID)
	}

	// admin topics disabled ?
	if adminTopicsEnabled == false {
		return
	}

	switch {

	// list all users
	case topic == "$SYS/auth/users/list/get":
		handled = true
		go onAuthUsersListGet(originClientID)

	case client.MatchRoute("$SYS/auth/user/#", topic):
		topicArray := strings.Split(topic, "/")
		username := topicArray[3]

		switch {

		// get a user to edit-it
		case client.MatchRoute("$SYS/auth/user/+/get", topic):
			handled = true
			go onAuthUserGet(originClientID, username)

		// set a user
		case client.MatchRoute("$SYS/auth/user/+/set", topic):
			handled = true
			go onAuthUserSet(originClientID, username, payload)

		// delete a user
		case client.MatchRoute("$SYS/auth/user/+/delete", topic):
			handled = true
			go onAuthUserDelete(originClientID, username)

		}

	}

	return
}
