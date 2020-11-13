package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// passwordCheck
func passwordCheck(username string, password string) (isOkay bool) {

	base64Data, exist := config.BcryptedPassword[username]
	if exist == true {
		basswordBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err == nil {
			errCompare := bcrypt.CompareHashAndPassword(basswordBytes, []byte(password))
			if errCompare == nil {
				isOkay = true
			}
		}
	}

	return
}
