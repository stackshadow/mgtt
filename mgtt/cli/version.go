package cli

import (
	"github.com/rs/zerolog/log"
)

// VersionFlag represents the flag which print the version
type VersionFlag bool

// AfterApply setup the json logging
func (v *VersionFlag) AfterApply() error {
	if bool(*v) == true {
		log.Info().Str("version", version).Send()
	}
	return nil
}
