package auth

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func configSave(filenameToLoad string) (err error) {
	var confidData []byte
	confidData, err = yaml.Marshal(config)
	if err == nil {
		if err = ioutil.WriteFile(filename, confidData, 0664); err != nil {
			log.Error().Str("filename", filename).Err(err).Msg("Error creating file")
		}
	}

	return
}
