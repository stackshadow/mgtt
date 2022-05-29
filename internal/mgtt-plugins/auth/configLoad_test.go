package auth

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	testConfigFile = `# Auth-plugin config-file

# uncomment this to enable anonym-login
# anonym: true

# use this to create a new user
new:
  - username: first
    password: firstsecret
    groups:
      - auth
      - debugging
  - username: second
    password: secondsecret
    groups:
      - debugging
`
)

func TestConfigLoad(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	Init()
	configLoad([]byte(testConfigFile))

	// check of config exist
	if isOkay := configCheckPassword("first", "firstsecret"); isOkay == false {
		t.FailNow()
	}
	if isOkay := configCheckPassword("second", "secondsecret"); isOkay == false {
		t.FailNow()
	}
}
