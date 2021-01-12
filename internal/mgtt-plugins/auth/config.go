package auth

import (
	"encoding/base64"
	"sync"

	"golang.org/x/crypto/bcrypt"
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
	New    []pluginConfigUser          `yaml:"new,omitempty"`
	Anonym bool                        `yaml:"anonym"`
	Users  map[string]pluginConfigUser `yaml:"users"`
}

type pluginConfigUser struct {
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password"`
	Groups   []string `yaml:"groups"`
}

// PasswordSet will conver tthe password-field to bcrypted-password
func (user *pluginConfigUser) PasswordSet(newPassword string) (err error) {
	var bcryptedData []byte
	bcryptedData, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = base64.StdEncoding.EncodeToString(bcryptedData)
	return
}

var mutex sync.Mutex
var filename string
var config *pluginConfig = &pluginConfig{
	Users: make(map[string]pluginConfigUser),
}
