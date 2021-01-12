package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onAuthUserSet(originClientID string, payload []byte) {

	var err error
	var newUserInfo pluginConfigUser
	var newUserPassword *string
	var newUserGroups *[]string

	err = json.Unmarshal(payload, &newUserInfo)

	if err == nil && newUserInfo.Username == "" {
		err = errors.New("Empty username")
	}

	// password
	if err == nil && newUserInfo.Password != "" {
		newUserPassword = &newUserInfo.Password
	}

	// save groups
	if err == nil && len(newUserInfo.Groups) > 0 {
		newUserGroups = &newUserInfo.Groups
	}

	// save
	if err == nil {
		userSet(newUserInfo.Username, newUserPassword, newUserGroups)
		err = configSave(filename)
	}

	if err == nil {

		err = clientlist.PublishToClient(
			originClientID,
			fmt.Sprintf("$SYS/auth/user/%s/set/success", newUserInfo.Username),
			[]byte("true"))
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
