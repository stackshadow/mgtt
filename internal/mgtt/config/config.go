package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// Config represents the config of your broker
var Values struct {
	URL string `yaml:"url" default:"tcp://0.0.0.0:8883"`

	TLS struct {
		CA   string `yaml:"ca" default:""` // path to ca-file for mTLS
		Cert string `yaml:"cert" default:"./mgtt.cert"`
		Key  string `yaml:"key" default:"./mgtt.key"`
	} `yaml:"tls"`

	DB string `yaml:"db" default:"./messages.db"`
}

var defaultConfig string = `
# The serve-url in the scheme tcp://<ip>:<port>
# as <ip> you usual will use 127.0.0.1 or 0.0.0.0
# as <port> you usual will use 8883
url: "tcp://0.0.0.0:8883"

tls:
  
  # if provided, mgtt use mTLS
  # if file not exist an CA will be created
  ca: ""

  # this is needed if you would like to use tls
  cert: "./mgtt.cert"

  # the private key, needed for tls
  key: "./mgtt.key"

# the db where to store persistant data
# this is needed for mqtt-persistand messages
db: "./messages.db"
`

// Load will load a file to Values
//
// if the file not exist, we save a defaultConfig with comments to <file>
func Load(file string) {

	var err error

	if file != "" {

		_, err = os.Stat(file)
		fileExist := errors.Is(err, os.ErrNotExist)

		var data []byte

		// read the file
		if fileExist {
			data, err = ioutil.ReadFile(file)
			utils.PanicOnErr(err)
		} else {
			data = []byte(defaultConfig)
			err = ioutil.WriteFile(file, data, 0600)
			utils.PanicOnErr(err)
		}

		// parse it
		err = yaml.Unmarshal(data, &Values)
		utils.PanicOnErr(err)

		// apply defaults
		defaults.SetDefaults(&Values)

	} else {
		log.Info().Msg("No filename provided, not loading config")
	}

}
