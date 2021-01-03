package auth

import (
	"fmt"
	"strings"

	"gitlab.com/mgtt/internal/mgtt/broker"
)

func onHandleUserDelete(originClientID string, topic string) {
	topicArray := strings.Split(topic, "/")
	username := topicArray[3]

	delete(config.BcryptedPassword, username)
	err := configSave(filename)

	if err == nil {
		if broker.Current != nil {
			broker.Current.PublishToClient(
				originClientID,
				fmt.Sprintf("$SYS/auth/user/%s/delete/ok", username),
				[]byte("true"))
		}
	}
}
