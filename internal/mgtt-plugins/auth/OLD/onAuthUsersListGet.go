package auth

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

type userListElement struct {
	Username string `json:"username"`
}

func onAuthUsersListGet(originClientID string) {
	var err error
	// a new list
	var newUserList []pluginConfigUser

	// fill the list
	for username, userinfo := range config.Users {
		userinfo.Username = username
		userinfo.Password = ""
		newUserList = append(newUserList, userinfo)
	}

	// create a json and send it
	jsonData, err := json.Marshal(newUserList)
	if err == nil {
		err = clientlist.PublishToClient(originClientID, "$SYS/auth/users/list/json", jsonData)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
