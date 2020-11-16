package auth

import (
	"sync"
)

const (
	defaultConfigContent = `# Auth-plugin config-file

# uncomment this to enable anonym-login
# anonym: true

# use this to create a new user
#new:
#  - username:
#    password:

`
)

type pluginConfig struct {
	New              []pluginConfigNewUser `yaml:"new,omitempty"`
	BcryptedPassword map[string]string     `yaml:"bcryptedpassword"`
	Anonym           bool                  `yaml:"anonym"`
}

type pluginConfigNewUser struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

var mutex sync.Mutex
var filename string
var config *pluginConfig = &pluginConfig{
	BcryptedPassword: make(map[string]string),
}
