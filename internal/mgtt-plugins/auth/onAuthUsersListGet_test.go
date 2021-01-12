package auth

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestOnAuthUsersListGet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Remove("./TestOnAuthUsersListGet_auth.yml")

	LocalInit("TestOnAuthUsersListGet_")

	// delete user
	OnHandleMessage("integrationtest", "$SYS/auth/users/list", []byte(""))

	os.Remove("./TestOnAuthUsersListGet_auth.yml")
}
