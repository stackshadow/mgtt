package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config represents the config of your broker
var Values struct {
	Level string `yaml:"level" default:"info"`
	JSON  bool   `yaml:"json" default:"false"`

	URL string `yaml:"url" default:"tcp://0.0.0.0:8883"`

	Timeout time.Duration `yaml:"timeout" default:"15s"`
	Retry   time.Duration `yaml:"retry" default:"30s"`

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

	Plugins map[string]interface{} `yaml:"plugins"`
}

// the raw configData
var valuesRawMap map[interface{}]interface{}

var fileName string

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
