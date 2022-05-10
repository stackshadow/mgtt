package config

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// Config represents the config of your broker
var Values struct {
	Level string `yaml:"level" default:"warn"`
	JSON  bool   `yaml:"json" default:"false"`

	URL string `yaml:"url" default:"tcp://0.0.0.0:8883"`

	Timeout time.Duration `yaml:"timeout" default:"15s"`
	Retry   time.Duration `yaml:"retry" default:"30s"`
	Plugins string        `yaml:"plugins" default:"auth,acl"`

	AdminTopics bool `yaml:"adminTopics" default:"false"`

	TLS struct {
		CA struct {
			File string `yaml:"file" default:""`

			Organization  string `yaml:"org" default:"FeelGood Inc."`
			Country       string `yaml:"country" default:"DE"`
			Province      string `yaml:"province" default:"Local"`
			Locality      string `yaml:"city" default:"Berlin"`
			StreetAddress string `yaml:"address" default:"Corner 42"`
			PostalCode    string `yaml:"code" default:"030423"`
		} `yaml:"ca"`

		Cert struct {
			File string `yaml:"file" default:""`

			Organization  string `yaml:"org" default:"FeelGood Inc."`
			Country       string `yaml:"country" default:"DE"`
			Province      string `yaml:"province" default:"Local"`
			Locality      string `yaml:"city" default:"Berlin"`
			StreetAddress string `yaml:"address" default:"Corner 42"`
			PostalCode    string `yaml:"code" default:"030423"`
		} `yaml:"cert"`
	} `yaml:"tls"`

	DB string `yaml:"db" default:"./messages.db"`
}

var defaultConfig string = `
# The serve-url in the scheme tcp://<ip>:<port>
# as <ip> you usual will use 127.0.0.1 or 0.0.0.0
# as <port> you usual will use 8883
url: "tcp://0.0.0.0:8883"

# Connection timeout for clients
timeout: 15

tls:
  
  # if provided, mgtt use mTLS
  # if file not exist an CA will be created
  ca:
    file: ""

  # this is needed if you would like to use tls
  cert:
    file: "./tls_cert.crt"

# the db where to store persistant data
# this is needed for mqtt-persistand messages
db: "./messages.db"

`

// Load will load a file to Values
//
// if the file not exist, we save a defaultConfig with comments to <file>
func MustLoad(file string) {

	var err error

	if file != "" {

		_, err = os.Stat(file)
		fileExist := !errors.Is(err, os.ErrNotExist)

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

	} else {
		log.Info().Msg("No filename provided, not loading config")
	}

}

// Apply log-level and log-output
func Apply() {

	var err error
	var newLogLevel zerolog.Level

	// apply defaults
	defaults.SetDefaults(&Values)

	// loglevel
	newLogLevel, err = zerolog.ParseLevel(Values.Level)
	if err == nil {
		zerolog.SetGlobalLevel(newLogLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	// jsonlog
	if !Values.JSON {
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

}
