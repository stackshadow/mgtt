package auth

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// userSet will set passwords/group for an user and return the new user-Object
//
// if the user already exist, we override the password
func userSet(username string, password *string, groups *[]string) (user pluginConfigUser, err error) {

	// get user
	user = config.Users[username]

	// we remove the username if it exist in the object
	user.Username = ""

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
