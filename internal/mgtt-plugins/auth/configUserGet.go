package auth

// configUserGet return the user without password
func configUserGet(username string) (user pluginConfigUser, exist bool) {
	if user, exist = config.Users[username]; exist {
		user.Username = username
		user.Password = ""
	}
	return
}
