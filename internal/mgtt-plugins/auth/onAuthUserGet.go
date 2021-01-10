package auth

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

type userElement struct {
	Username string   `json:"username"`
	Password string   `json:"password,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

func onAuthUserGet(originClientID string, username string) {
	var err error

	// check if the user exist
	if _, exist := config.BcryptedPassword[username]; exist {

		var user = userElement{
			Username: username,
			Groups:   config.groups[username],
		}

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
