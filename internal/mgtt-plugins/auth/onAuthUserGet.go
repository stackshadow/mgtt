package auth

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onAuthUserGet(originClientID string, username string) {
	var err error

	// check if the user exist
	if user, exist := config.Users[username]; exist {

		// remove the password, nobody should now about it
		user.Username = username
		user.Password = ""

		// create a json and send it
		var jsonData []byte
		jsonData, err = json.Marshal(user)
		if err == nil {
			err = clientlist.PublishToClient(
				originClientID,
				fmt.Sprintf("$SYS/auth/user/%s/json", username),
				jsonData,
			)
		}

	} else {
		err = clientlist.PublishToClient(
			originClientID,
			fmt.Sprintf("$SYS/auth/user/%s/error", username),
			[]byte("User dont exist"),
		)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
