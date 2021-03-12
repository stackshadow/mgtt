package broker

import (
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func TestLastWill(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestLastWill_test.db")
	defer os.Remove("TestLastWill_test.db")
	server, _ := New()
	serverConfig := Config{
		URL:        "tcp://127.0.0.1:1234",
		DBFilename: "TestLastWill_test.db",
	}
	server.Serve(serverConfig)

	// ###############################################  The client with will-message
	var pahoClientConnected bool = false

	pahoClient, err := testConnectClient(serverConfig.URL, "", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// ###############################################  The client
	pahoClientSub, err := testConnectClient(serverConfig.URL, "", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var lastWillRecvd sync.Mutex
	lastWillRecvd.Lock()
	if token := pahoClientSub.Subscribe("lastwill/#", 0, func(client paho.Client, msg paho.Message) {
		if string(msg.Payload()) != "value" {
			t.FailNow()
			return
		}
		lastWillRecvd.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	// disconnect first client and wait for last-will
	pahoClient.Disconnect(500)
	lastWillRecvd.Lock()

	// disconnect the last client
	pahoClientSub.Disconnect(500)

	// close the server
	server.ServeClose()
}
