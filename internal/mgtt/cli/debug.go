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
	}
	return nil
}
