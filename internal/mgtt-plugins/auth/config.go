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
#  - username: johndoe
#    password: mysupersecretpassword
#    groups:
#      - auth
#      - debugging

`
)

type pluginConfig struct {
	New    []pluginConfigUser          `yaml:"new,omitempty"`
	Anonym bool                        `yaml:"anonym"`
	Users  map[string]pluginConfigUser `yaml:"users"`
}

type pluginConfigUser struct {
	Username string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password" json:"password,omitempty"`
	Groups   []string `yaml:"groups" json:"groups"`
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
