package auth

import (
	"strings"

	"gitlab.com/mgtt/internal/mgtt/client"
)

// OnHandleMessage gets called after OnPublishRequest
//
// If this function return true, the plugin handled the message and no other plugin will get it
//
// If a plugin handle the message, it will NOT sended to subscribers
func OnHandleMessage(originClientID string, topic string, payload []byte) (handled bool) {

	switch {

	// who is currently logged in
	case topic == "$SYS/self/username/get":
		// topic matched, we handled it
		handled = true
		go onSelfUsernameGet(originClientID)

		// who is currently logged in
	case topic == "$SYS/self/groups/get":
		handled = true
		go onSelfGroupsGet(originClientID)

	// list all users
	case topic == "$SYS/auth/users/list/get":
		// topic matched, we handled it
		handled = true
		go onHandleUserList(originClientID)

	// get a user to edit-it
	case client.MatchRoute("$SYS/auth/user/+/get", topic):
		topicArray := strings.Split(topic, "/")
		username := topicArray[3]
		handled = true
		go onHandleUserGet(originClientID, username)

	// set a new password
	case client.MatchRoute("$SYS/auth/user/+/password/set", topic):
		// topic matched, we handled it
		handled = true
		go onHandlePasswordSet(originClientID, topic, string(payload))

	// delete a user
	case client.MatchRoute("$SYS/auth/user/+/delete", topic):
		// topic matched, we handled it
		handled = true
		go onHandleUserDelete(originClientID, topic)

	}

	return
}
