package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// passwordAdd will add a new user with an password
//
// if the user already exist, we override the password
func passwordAdd(username string, password string) (err error) {

	// convert passwort to base64-bcrypt
	var bcryptedData []byte
	bcryptedData, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	bcryptedBase64 := base64.StdEncoding.EncodeToString(bcryptedData)

	// save it to the config
	config.BcryptedPassword[username] = bcryptedBase64

	return
}
