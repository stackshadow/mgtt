package client

// UsernameSet set the username
//
// username can only be set, if its not set before !
func (c *MgttClient) UsernameSet(username string) {
	if c.username == "" {
		c.username = username
	}
}

// Username return the current username
func (c *MgttClient) Username() string {
	return c.username
}
