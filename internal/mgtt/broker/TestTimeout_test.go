package broker

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestTimeout(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// setup connection timeout
	ConnectTimeout = 1

	// the broker
	os.Remove("TestTimeout_test.db")
	defer os.Remove("TestTimeout_test.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1240",
			DBFilename: "TestTimeout_test.db",
		},
	)
	time.Sleep(time.Second * 2)

	connection, err := net.Dial("tcp", "127.0.0.1:1240")
	if err != nil {
		t.FailNow()
	}

	// now we are connected
	time.Sleep(time.Second * 2)
	msg := make([]byte, 4000)
	_, err = connection.Read(msg)
	if err == nil {
		t.FailNow()
	}

	server.ServeClose()
	time.Sleep(time.Second * 1)
}
