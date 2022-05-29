package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// MustLoadFromFile will load from <file> if not set to "" and panics on an error
func MustLoadFromFile(file string) {

	var err error

	if file != "" {

		_, err = os.Stat(file)
		fileExist := !errors.Is(err, os.ErrNotExist)

		var data []byte

		// read the file
		if fileExist {
			data, err = ioutil.ReadFile(file) //#nosec
			utils.PanicOnErr(err)
		}

		MustLoadFromString(string(data)) // load from readed bytes
		Globals.fileName = file
	} else {
		log.Info().Msg("No filename provided, not loading config from a file")
	}

}

// MustLoadFromString will load the config from a json-string and panics if an error occurred
func MustLoadFromString(data string) {
	if data == "" {
		return
	}

	// parse it to
	err := yaml.Unmarshal([]byte(data), &Globals)
	utils.PanicOnErr(err)
}
