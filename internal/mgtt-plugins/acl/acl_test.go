package acl

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestPublish(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	Init()

	// we add some acls
	pluginConfig.Rules["testuser"] = []aclEntry{
		// we not allow write to clients
		{
			Direction: "w",
			Route:     "$SYS/broker/clients",
			Allow:     false,
		},

		// and not to the rest
		{
			Direction: "w",
			Route:     "$SYS/#",
			Allow:     false,
		},

		// but we allow write to sensors
		{
			Direction: "w",
			Route:     "sensors/#",
			Allow:     true,
		},

		// and not to the rest
		{
			Direction: "w",
			Route:     "#",
			Allow:     false,
		},
	}

	if OnPublishRequest("", "testuser", "$SYS/broker/clients") == true {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "$SYS/connections/count") == true {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "sensors/temp/first") == false {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "sensors/temp/second") == false {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "commands/power/off") == true {
		t.FailNow()
	}

	// we add some acls
	pluginConfig.Rules["_anonym"] = []aclEntry{
		// we not allow write to clients
		{
			Direction: "r",
			Route:     "system/users",
			Allow:     true,
		},
		{
			Direction: "w",
			Route:     "system/config",
			Allow:     true,
		},
	}

	if OnPublishRequest("", "", "system/users") == true {
		t.FailNow()
	}
	if OnSendToSubscriberRequest("", "", "system/users") == false {
		t.FailNow()
	}

	if OnPublishRequest("", "", "system/config") == false {
		t.FailNow()
	}
	if OnSendToSubscriberRequest("", "", "system/config") == true {
		t.FailNow()
	}

}
