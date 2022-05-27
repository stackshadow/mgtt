package persistance

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestClientSession(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.Logger.With().Caller().Logger()

	// delete db if exist
	_, err := os.Stat("integrtest.db")
	if err == nil {
		os.Remove("integrtest.db")
	}
	defer os.Remove("integrtest.db")

	// open DB
	MustOpen("integrtest.db")

	// set subscriptions
	SubscriptionsSet("id1", []string{"/a", "/b"})
	SubscriptionsSet("id2", []string{"/c", "/d", "/e"})

	if subscriptions := SubscriptionsGet("notexist"); len(subscriptions) != 0 {
		t.FailNow()
	}

	if subscriptions := SubscriptionsGet("id1"); len(subscriptions) != 2 {
		t.FailNow()
	}

	if subscriptions := SubscriptionsGet("id2"); len(subscriptions) != 3 {
		t.FailNow()
	}

	// clean
	CleanSession("id1")

	if subscriptions := SubscriptionsGet("id1"); len(subscriptions) != 0 {
		t.FailNow()
	}
	if subscriptions := SubscriptionsGet("id2"); len(subscriptions) != 3 {
		t.FailNow()
	}
}
