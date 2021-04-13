package cli

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DebugFlag represents the flag which enable debugging
type DebugFlag bool

// AfterApply setup the json logging
func (v *DebugFlag) AfterApply() error {

	// we add the caller by default
	log.Logger = log.Logger.With().Caller().Logger()

	if bool(*v) == true {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Enable log-level: debug")
		/*
			log.Logger = log.Logger.Output(&lumberjack.Logger{
				Filename:   "./foo.log",
				MaxSize:    1,  // megabytes after which new file is created
				MaxBackups: 3,  // number of backups
				MaxAge:     28, //days
			})
		*/

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("Enable log-level: info")
	}
	return nil
}
