package auth

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// configLoad will load an file
func configLoad(fileData []byte) (changed bool) {

	mutex.Lock()
	defer mutex.Unlock()

	err := yaml.Unmarshal(fileData, &pluginConfig)
	utils.PanicOnErr(err)

	// add new users
	for _, newUser := range pluginConfig.New {
		if newUser.Username != "" && newUser.Password != "" {
			log.Debug().Str("username", newUser.Username).Msg("found new user, store it to config")
			newUsername := newUser.Username
			_, err = configSetUser(newUsername, &newUser.Password, &newUser.Groups)
			changed = true
		}
	}
	pluginConfig.New = []pluginConfigUser{}

	log.Info().Msg("Loaded config")
	return
}
