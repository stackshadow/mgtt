package auth

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// configLoad will load an file
func configLoad(fileData []byte) {

	mutex.Lock()
	defer mutex.Unlock()

	err := yaml.Unmarshal(fileData, config)
	utils.PanicOnErr(err)

	// add new users
	for _, newUser := range config.Plugins.ACL.New {
		if newUser.Username != "" && newUser.Password != "" {
			newUsername := newUser.Username
			_, err = configSetUser(newUsername, &newUser.Password, &newUser.Groups)
		}
	}
	config.Plugins.ACL.New = []pluginConfigUser{}

	log.Info().Msg("Loaded config")
}
