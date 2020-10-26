package auth

import (
	"encoding/base64"
	"fmt"
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

	filename string
	mutex    sync.Mutex
}

var config *pluginConfig = &pluginConfig{
	BcryptedPassword: make(map[string]string),
}

// loadConfig will load an file
func loadConfig(filename string) (err error) {

	config.mutex.Lock()
	defer config.mutex.Unlock()

	// store filename
	config.filename = filename

	var fileData []byte
	fileData, err = ioutil.ReadFile(config.filename)
	if err != nil {
		log.Warn().Err(err).Msg("Error opening config file")
		log.Info().Msg("Creating default auth.yml file")
		if err = ioutil.WriteFile(config.filename, []byte(defaultConfigContent), 0664); err != nil {
			log.Warn().Err(err).Msg("Error creating default file")
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

	return
}

// watchConfig will watch for changes of the config file and reload it when it changes
func watchConfig() (err error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	w.Add(config.filename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				loadConfig(config.filename)

			case err := <-w.Error:
				log.Err(err).Send()
			case <-w.Closed:
				return
			}
		}
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Err(err).Send()
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
		if err = ioutil.WriteFile(config.filename, confidData, 0664); err != nil {
			log.Err(err).Msg("Error creating default auth.yml file")
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
