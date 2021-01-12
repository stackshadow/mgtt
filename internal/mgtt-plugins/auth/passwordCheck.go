package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// passwordCheck
func passwordCheck(username string, password string) (isOkay bool) {

	if config.Anonym == true && username == "" {
		return true
	}

	// get user
	var user = config.Users[username]

	basswordBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err == nil {
		errCompare := bcrypt.CompareHashAndPassword(basswordBytes, []byte(password))
		if errCompare == nil {
			isOkay = true
		}
	}

	return
}
