package auth

import (
	"fmt"
	"strings"

	"gitlab.com/mgtt/broker"
)

func onHandleUserDelete(originClientID string, topic string ) {
	topicArray := strings.Split(topic, "/")
	username := topicArray[3]

	delete(config.BcryptedPassword, username)
	err := saveConfig(filename)

	if err == nil {
		if broker.Current != nil {
			broker.Current.PublishToClient(
				originClientID,
				fmt.Sprintf("$SYS/auth/user/%s/delete/ok", username),
				[]byte("true"))
		}
	}
}
