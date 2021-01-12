package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// passwordAdd will add a new user with an password
//
// if the user already exist, we override the password
func userSet(username string, password *string, groups *[]string) (err error) {

	// get or create user
	var user pluginConfigUser

	// get user
	user = config.Users[username]

	// password
	if password != nil {
		var bcryptedData []byte
		bcryptedData, err = bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)

		user.Password = base64.StdEncoding.EncodeToString(bcryptedData)
	}

	// groups
	if groups != nil {
		user.Groups = *groups
	}

	// save it to the config
	config.Users[username] = user

	return
}
