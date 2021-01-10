package auth

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onSelfGroupsGet(originClientID string) {
	var err error
	userName := broker.Current.UserNameOfClient(originClientID)

	userGroups, exist := config.groups[userName]
	if exist == false {
		userGroups = append(userGroups, "anonym")
	}

	var userGroupsBytes []byte
	userGroupsBytes, err = json.Marshal(userGroups)
	if err != nil {
		err = clientlist.PublishToClient(
			originClientID,
			"$SYS/self/groups/json",
			userGroupsBytes)
	}

	if err != nil {
		log.Error().Err(err).Send()
	}
}
