package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

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
			data, err = ioutil.ReadFile(file) //#nosec
			utils.PanicOnErr(err)
		}
		fileName = file

		// load from readed bytes
		LoadFromByte(data)

	} else {
		log.Info().Msg("No filename provided, not loading config")
	}

}

func LoadFromByte(data []byte) {
	var err error

	// parse it to raw values
	err = yaml.Unmarshal(data, &valuesRawMap)
	utils.PanicOnErr(err)

	// parse to Values
	err = yaml.Unmarshal(data, &Values)
	utils.PanicOnErr(err)

	// apply defaults
	defaults.SetDefaults(&Values)

	// setup logs
	ApplyLog()

	// call plugins
	plugin.CallOnPluginConfig(Values.Plugins)
	MustSave() // if config was changed

}
