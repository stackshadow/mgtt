package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// configCheckPassword
func configCheckPassword(username string, password string) (isOkay bool) {

	if config.Plugins.ACL.Anonym == true && username == "" {
		return true
	}

	// get user
	var user = config.Plugins.ACL.Users[username]

	basswordBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err == nil {
		errCompare := bcrypt.CompareHashAndPassword(basswordBytes, []byte(password))
		if errCompare == nil {
			isOkay = true
		}
	}

	return
}
