package acl

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigContent = `# ACL-plugin config-file

# This is an example for an acl
rules:

# The username
  admin:

# here we define an route, according to the mqtt-schema with # and + can be used
	- route: "#"
	
# w means the publisher write to the broker
# r means the broker send the message to an subscriber
	  direction: w

# and here we decide what to do if the route matches
	  allow: true

    - route: "#"
      direction: r
      allow: true

# This is an anonymouse user
#  _anonym:
#	- route: "#"
#      direction: w
#      allow: true

#	- route: "#"
#      direction: r
#      allow: true


`
)

type pluginConfig struct {
	Rules map[string][]aclEntry `yaml:"rules"`
}

var mutex sync.Mutex
var filename string
var config *pluginConfig = &pluginConfig{
	Rules: make(map[string][]aclEntry),
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
