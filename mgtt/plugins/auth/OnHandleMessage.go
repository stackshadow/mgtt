package auth

import (
	"gitlab.com/mgtt/broker"
	"gitlab.com/mgtt/client"
)

type userListElement struct {
	Username string `json:"username"`
}

// OnHandleMessage gets called after OnPublishRequest
//
// If this function return true, the plugin handled the message and no other plugin will get it
//
// If a plugin handle the message, it will NOT sended to subscribers
func OnHandleMessage(originClientID string, topic string, payload []byte) (handled bool) {

	switch {

	// who is currently logged in
	case topic == "$SYS/auth/user/whoami":
		// topic matched, we handled it
		handled = true
		broker.Current.PublishToClient(
			originClientID,
			"$SYS/auth/user/name",
			[]byte(broker.Current.UserNameOfClient(originClientID)),
		)

	// list all users
	case topic == "$SYS/auth/users/list":
		// topic matched, we handled it
		handled = true
		go onHandleUserList(originClientID)

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
