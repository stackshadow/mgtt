package auth

import (
	"encoding/json"

	"gitlab.com/mgtt/broker"
)

func onHandleUserList(originClientID string) {
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
			broker.Current.PublishToClient(originClientID, "$SYS/auth/users/list/json", jsonData)
		}
	}
}
