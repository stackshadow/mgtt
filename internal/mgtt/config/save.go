package config

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

func MustSave() {

	if Globals.fileName != "" {

		var err error
		var data []byte

		// parse it
		data, err = yaml.Marshal(Globals)
		utils.PanicOnErr(err)

		// save it to file
		err = ioutil.WriteFile(Globals.fileName, data, 0600)
		utils.PanicOnErr(err)

		log.Info().
			Str("filename", Globals.fileName).
			Msg("Config file changes, save it")
	}
}
