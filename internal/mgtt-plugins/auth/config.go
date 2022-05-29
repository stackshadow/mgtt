package auth

import (
	"encoding/base64"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type pluginConfigStruct struct {
	Enable bool                        `yaml:"enable"`
	New    []pluginConfigUser          `yaml:"new,omitempty"`
	Anonym bool                        `yaml:"anonym"`
	Users  map[string]pluginConfigUser `yaml:"users,omitempty"`
}

type pluginConfigUser struct {
	Username string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password" json:"password,omitempty"`
	Groups   []string `yaml:"groups" json:"groups"`
}

// PasswordSet will convert the password-field to bcrypted-password
func (user *pluginConfigUser) PasswordSet(newPassword string) (err error) {
	var bcryptedData []byte
	bcryptedData, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = base64.StdEncoding.EncodeToString(bcryptedData)
	return
}

var mutex sync.Mutex
var pluginConfig pluginConfigStruct = pluginConfigStruct{}
