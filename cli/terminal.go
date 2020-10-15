package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TerminalFlag represents the flag if we log to terminal or as json
type TerminalFlag bool

// AfterApply setup the json logging
func (v *TerminalFlag) AfterApply() error {
	if bool(*v) == true {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return nil
}
