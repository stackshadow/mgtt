package config

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

func MustSave() {
	if fileName != "" {

		var err error
		var data []byte

		// values to data
		data, err = yaml.Marshal(Values)
		utils.PanicOnErr(err)

		// data to map
		err = yaml.Unmarshal(data, &valuesRawMap)
		utils.PanicOnErr(err)

		// parse it
		data, err = yaml.Marshal(&valuesRawMap)
		utils.PanicOnErr(err)

		// save it to file
		err = ioutil.WriteFile(fileName, data, 0600)
		utils.PanicOnErr(err)

		log.Info().
			Str("filename", fileName).
			Msg("Config file changes, save it")
	}
}
