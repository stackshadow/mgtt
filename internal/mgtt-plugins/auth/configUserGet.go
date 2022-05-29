package auth

// configUserGet return the user without password
func configUserGet(username string) (user pluginConfigUser, exist bool) {

	if username == "" {
		user.Username = "_anonym"
		user.Password = ""

		exist = true
		return
	}

	if user, exist = pluginConfig.Users[username]; exist {
		user.Username = username
		user.Password = ""
	}
	return
}
