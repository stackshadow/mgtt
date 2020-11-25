package broker

import (
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func TestTimeout(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// create a plugin that unlocks the connection of a client
	var connectedLock sync.Mutex
	var connectedLockPlugin plugin.V1
	connectedLockPlugin.OnNewClient = func(remoteAddr string) {
		connectedLock.Unlock()
		return
	}
	plugin.Register("connectedLockPlugin", &connectedLockPlugin)

	// setup connection timeout
	ConnectTimeout = 1

	// the broker
	os.Remove("test1.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1240",
			DBFilename: "test1.db",
		},
	)
	time.Sleep(time.Second * 2)

	connectedLock.Lock()
	connection, err := net.Dial("tcp", "127.0.0.1:1240")
	if err != nil {
		t.FailNow()
	}
	connectedLock.Lock()

	// now we are connected
	msg := make([]byte, 4000)
	_, err = connection.Read(msg)
	if err == nil {
		t.FailNow()
	}

}
