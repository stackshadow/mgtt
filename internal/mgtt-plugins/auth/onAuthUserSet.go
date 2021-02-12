package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onAuthUserSet(originClientID string, username string, payload []byte) {

	var err error
	var newUserInfo pluginConfigUser
	var newUserPassword *string
	var newUserGroups *[]string

	err = json.Unmarshal(payload, &newUserInfo)

	// username
	if err == nil && username == "" {
		err = errors.New("Empty username")
	}

	// password
	if err == nil && newUserInfo.Password != "" {
		newUserPassword = &newUserInfo.Password
	}

	// save groups
	if err == nil {
		newUserGroups = &newUserInfo.Groups
	}

	// save
	var changedUser pluginConfigUser
	if err == nil {
		changedUser, _ = userSet(username, newUserPassword, newUserGroups)
		err = configSave(filename)
	}

	if err == nil {

		changedUser.Username = username

		// create a json and send it
		var jsonData []byte
		jsonData, err = json.Marshal(changedUser)

		err = clientlist.PublishToClient(
			originClientID,
			fmt.Sprintf("$SYS/auth/user/%s/set/success", username),
			jsonData)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
