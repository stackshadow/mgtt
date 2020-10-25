package auth

import (
	"io/ioutil"
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

	os.Remove(config.filename)
}

const (
	newconfig = `# Auth-plugin config-file

# use this to create a new user
new:
  username: testuser
  password: dummypassword

`
)

func TestChangeOfConfigFile(t *testing.T) {

	loadConfig("integrationtest.yml")

	// add a new user
	if err := ioutil.WriteFile(config.filename, []byte(newconfig), 0664); err != nil {
		t.FailNow()
	}
	loadConfig(config.filename)

	if passwordCheck("testuser", "dummypassword") == false {
		t.FailNow()
	}

	os.Remove(config.filename)
}
