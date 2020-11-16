package auth

import (
	"fmt"
	"strings"

	"gitlab.com/mgtt/broker"
)

func onHandlePasswordSet(originClientID string, topic string, payload string) {
	topicArray := strings.Split(topic, "/")
	username := topicArray[3]

	err := passwordAdd(username, payload)
	if err == nil {
		err = saveConfig(filename)
	}

	if err == nil {
		if broker.Current != nil {
			broker.Current.PublishToClient(
				originClientID,
				fmt.Sprintf("$SYS/auth/user/%s/password/set/ok", username),
				[]byte("true"))
		}
	}
}
