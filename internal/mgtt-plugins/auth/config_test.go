package auth

import (
	"os"
	"testing"
)

func TestConfigFile(t *testing.T) {

	configLoad("integrationtest.yml")

	// add a new user
	var userPassword = "testpassword"
	userSet("testuser", &userPassword, nil)

	if passwordCheck("testuser", userPassword) == false {
		t.FailNow()
	}

	os.Remove(filename)
}

func TestEnvironment(t *testing.T) {

	os.Setenv("AUTH_USERNAME", "envuser")
	os.Setenv("AUTH_PASSWORD", "envpw")
	os.Setenv("ENABLE_ADMIN_TOPICS", "true")

	OnInit("integrationtest.yml")

	if passwordCheck("envuser", "envpw") == false {
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
