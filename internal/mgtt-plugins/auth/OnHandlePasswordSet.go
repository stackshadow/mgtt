package auth

import (
	"fmt"

	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onHandlePasswordSet(originClientID string, topic string, payload string) {

	var err error

	topicArray := strings.Split(topic, "/")
	username := topicArray[3]

	err = passwordAdd(username, payload)
	if err == nil {
		err = configSave(filename)
	}

	if err == nil {

		err = clientlist.PublishToClient(
			originClientID,
			fmt.Sprintf("$SYS/auth/user/%s/password/set/success", username),
			[]byte("true"))
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
