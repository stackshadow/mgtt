package acl

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

// configLoad will load an file
func configLoad(fileData []byte) {

	mutex.Lock()
	defer mutex.Unlock()

	err := yaml.Unmarshal(fileData, config)
	utils.PanicOnErr(err)

	log.Info().Msg("Loaded config")
}
