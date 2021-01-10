package auth

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onAuthUserPasswordSet(originClientID string, username string, password string) {

	var err error

	err = passwordAdd(username, password)
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
