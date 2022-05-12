package auth

import (
	"testing"
)

func TestConfigFile(t *testing.T) {

	Init()

	// add a new user
	var userPassword = "testpassword"
	configSetUser("testuser", &userPassword, nil)

	if configCheckPassword("testuser", userPassword) == false {
		t.FailNow()
	}

}
