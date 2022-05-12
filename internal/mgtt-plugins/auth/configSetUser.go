package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// userSet will set passwords/group for an user and return the new user-Object
//
// if the user already exist, we override the password
func configSetUser(username string, password *string, groups *[]string) (user pluginConfigUser, err error) {

	// get user
	user, _ = configUserGet(username)

	// you can not set the anonymouse-password
	if username == "" {
		return
	}

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
	config.Plugins.ACL.Users[username] = user

	// get user
	user, _ = configUserGet(username)
	return
}
