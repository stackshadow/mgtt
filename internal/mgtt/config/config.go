package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt-plugins/acl"
	"gitlab.com/mgtt/internal/mgtt-plugins/auth"
	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// Config represents the config of your broker
var Values struct {
	Level string `yaml:"level" default:"info"`
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
		}

		// parse it
		err = yaml.Unmarshal(data, &Values)
		utils.PanicOnErr(err)

		// apply defaults
		defaults.SetDefaults(&Values)

		// setup logs
		ApplyLog()

		// plugins
		pluginList := strings.Split(Values.Plugins, ",")
		for _, pluginName := range pluginList {
			if pluginName == "acl" {
				acl.Init()
			}
			if pluginName == "auth" {
				auth.Init()
			}
		}

		// call plugins
		plugin.CallOnConfig(data)

	} else {
		log.Info().Msg("No filename provided, not loading config")
	}

}

// Apply log-level and log-output
func ApplyLog() {

	var err error
	var newLogLevel zerolog.Level

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
