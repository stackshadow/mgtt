package auth

import (
	"io/ioutil"
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

	os.Remove("./TestConfigLoad.yml")

	if err := ioutil.WriteFile("./TestConfigLoad.yml", []byte(testConfigFile), 0664); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := configLoad("./TestConfigLoad.yml"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// check of config exist
	if isOkay := passwordCheck("first", "firstsecret"); isOkay == false {
		t.FailNow()
	}
	if isOkay := passwordCheck("second", "secondsecret"); isOkay == false {
		t.FailNow()
	}
}
