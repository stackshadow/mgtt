package auth

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onSelfUserGet(originClientID string) {
	var err error
	var username = clientlist.Get(originClientID).Username()

	// check if the user exist
	if user, exist := configUserGet(username); exist {

		// create a json and send it
		var jsonData []byte
		jsonData, err = json.Marshal(user)
		if err == nil {
			err = clientlist.PublishToClient(
				originClientID,
				"$SYS/self/user/json",
				jsonData,
			)
		}

	} else {
		err = clientlist.PublishToClient(
			originClientID,
			"$SYS/self/user/error",
			[]byte("User dont exist"),
		)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
