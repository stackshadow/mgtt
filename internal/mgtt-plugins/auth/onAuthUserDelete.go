package auth

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onAuthUserDelete(originClientID string, username string) {
	var err error

	delete(config.Users, username)
	configSave(filename)

	err = clientlist.PublishToClient(
		originClientID,
		fmt.Sprintf("$SYS/auth/user/%s/delete/success", username),
		[]byte("true"))

	if err != nil {
		log.Error().Err(err).Send()
	}

}
