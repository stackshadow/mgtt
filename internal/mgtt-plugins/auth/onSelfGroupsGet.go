package auth

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onSelfGroupsGet(originClientID string) {
	var err error
	var userGroups []string
	var userName = broker.Current.UserNameOfClient(originClientID)

	// check if the user exist
	if user, exist := config.Users[userName]; exist == true {
		userGroups = user.Groups
	} else {
		userGroups = append(userGroups, "anonym")
	}

	var userGroupsBytes []byte
	userGroupsBytes, err = json.Marshal(userGroups)
	if err == nil {
		err = clientlist.PublishToClient(
			originClientID,
			"$SYS/self/groups/json",
			userGroupsBytes)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
