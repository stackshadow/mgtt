package cli

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DebugFlag represents the flag which enable debugging
type DebugFlag bool

// AfterApply setup the json logging
func (v *DebugFlag) AfterApply() error {
	if bool(*v) == true {
		log.Logger = log.Logger.With().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	return nil
}
