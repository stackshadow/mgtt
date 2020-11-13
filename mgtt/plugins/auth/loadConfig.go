package auth

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// loadConfig will load an file
func loadConfig(filenameToLoad string) (err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// store filename
	filename = filenameToLoad

	var fileData []byte
	fileData, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Warn().Str("filename", filenameToLoad).Err(err).Msg("Error opening config file")
		log.Info().Str("filename", filenameToLoad).Msg("Creating default file")

		saveConfig(filenameToLoad)
	}

	err = yaml.Unmarshal(fileData, config)
	if err == nil {

		newUserExist := false
		for _, newUser := range config.New {
			if newUser.Username != "" && newUser.Password != "" {
				newUsername := newUser.Username
				newPassword := newUser.Password
				err = passwordAdd(newUsername, newPassword)
				newUserExist = true
			}
		}
		config.New = []pluginConfigNewUser{}

		if newUserExist == true {
			saveConfig(filenameToLoad)
		}
	}

	log.Info().Str("filename", filenameToLoad).Msg("Loaded config file")
	return
}
