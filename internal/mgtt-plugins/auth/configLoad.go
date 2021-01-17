package auth

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// configLoad will load an file
func configLoad(filenameToLoad string) (err error) {

	mutex.Lock()
	defer mutex.Unlock()

	// remember the filename globally
	filename = filenameToLoad

	// read the file-content
	var fileData []byte
	fileData, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Warn().Str("filename", filenameToLoad).Err(err).Msg("Error opening config file")
		log.Info().Str("filename", filenameToLoad).Msg("Creating default file")
		if err = ioutil.WriteFile(filenameToLoad, []byte(defaultConfigContent), 0664); err != nil {
			log.Warn().Str("filename", filenameToLoad).Err(err).Msg("Error creating default file")
		}
		fileData = []byte(defaultConfigContent)
	}

	err = yaml.Unmarshal(fileData, config)
	if err == nil {

		// if this gets true, we save the config file
		newUserExist := false

		// add new users
		for _, newUser := range config.New {
			if newUser.Username != "" && newUser.Password != "" {
				newUsername := newUser.Username
				_, err = userSet(newUsername, &newUser.Password, &newUser.Groups)
				newUserExist = true
			}
		}
		config.New = []pluginConfigUser{}

		// save the changed config-file
		if newUserExist == true {
			configSave(filenameToLoad)
		}
	}

	log.Info().Str("filename", filenameToLoad).Msg("Loaded config file")
	return
}
