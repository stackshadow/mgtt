package auth

import (
	"encoding/base64"
	"io/ioutil"

	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigContent = `# Auth-plugin config-file

# use this to create a new user
#new:
#  username:
#  password:

`
)

type pluginConfig struct {
	New struct {
		Username string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
	} `yaml:"new,omitempty"`

	BcryptedPassword map[string]string
}

var mutex sync.Mutex
var filename string
var config *pluginConfig = &pluginConfig{
	BcryptedPassword: make(map[string]string),
}

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
		if err = ioutil.WriteFile(filename, []byte(defaultConfigContent), 0664); err != nil {
			log.Warn().Str("filename", filenameToLoad).Err(err).Msg("Error creating default file")
		}
		fileData = []byte(defaultConfigContent)
	}

	err = yaml.Unmarshal(fileData, config)
	if err == nil {
		if config.New.Username != "" && config.New.Password != "" {
			newUsername := config.New.Username
			newPassword := config.New.Password

			config.New.Username = ""
			config.New.Password = ""
			err = passwordAdd(newUsername, newPassword)
		}
	}

	log.Info().Str("filename", filenameToLoad).Msg("Loaded config file")
	return
}

// watchConfig will watch for changes of the config file and reload it when it changes
func watchConfig() (err error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	w.Add(filename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Info().Str("filename", event.Path).Str("event", event.String()).Msg("File change detected")
				loadConfig(filename)

			case err := <-w.Error:
				log.Error().Err(err).Send()
			case <-w.Closed:
				return
			}
		}
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Error().Err(err).Send()
	}

	return
}

// passwordAdd will add a new user with an password
func passwordAdd(username string, password string) (err error) {

	// convert passwort to base64-bcrypt
	var bcryptedData []byte
	bcryptedData, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	bcryptedBase64 := base64.StdEncoding.EncodeToString(bcryptedData)

	// save it to the config
	config.BcryptedPassword[username] = bcryptedBase64

	confidData, err := yaml.Marshal(config)
	if err == nil {
		if err = ioutil.WriteFile(filename, confidData, 0664); err != nil {
			log.Error().Str("filename", filename).Err(err).Msg("Error creating file")
		}
	}

	return
}

// passwordCheck
func passwordCheck(username string, password string) (isOkay bool) {

	base64Data, exist := config.BcryptedPassword[username]
	if exist == true {
		basswordBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err == nil {
			errCompare := bcrypt.CompareHashAndPassword(basswordBytes, []byte(password))
			if errCompare == nil {
				isOkay = true
			}
		}
	}

	return
}
