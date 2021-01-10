package auth

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onHandleUserDelete(originClientID string, topic string) {
	var err error

	topicArray := strings.Split(topic, "/")
	username := topicArray[3]

	delete(config.BcryptedPassword, username)
	configSave(filename)

	err = clientlist.PublishToClient(
		originClientID,
		fmt.Sprintf("$SYS/auth/user/%s/delete/ok", username),
		[]byte("true"))

	if err != nil {
		log.Error().Err(err).Send()
	}
}
