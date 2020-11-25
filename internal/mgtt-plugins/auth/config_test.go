package auth

import (
	"os"
	"testing"
)

func TestConfigFile(t *testing.T) {

	loadConfig("integrationtest.yml")

	// add a new user
	passwordAdd("testuser", "testpassword")

	if passwordCheck("testuser", "testpassword") == false {
		t.FailNow()
	}

	os.Remove(filename)
}

const (
	newconfig = `# Auth-plugin config-file

# use this to create a new user
new:
  username: testuser
  password: dummypassword

`
)
