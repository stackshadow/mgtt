package auth

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

type userListElement struct {
	Username string `json:"username"`
}

func onHandleUserList(originClientID string) {
	var err error
	// a new list
	var newUserList []userListElement

	// fill the list
	for username := range config.BcryptedPassword {
		newUserList = append(newUserList, userListElement{
			Username: username,
		})
	}

	// create a json and send it
	jsonData, err := json.Marshal(newUserList)
	if err == nil {
		if broker.Current != nil {
			err = clientlist.PublishToClient(originClientID, "$SYS/auth/users/list/json", jsonData)
		}
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
