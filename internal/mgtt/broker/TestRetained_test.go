package broker

import (
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func TestRetained(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestRetained_test.db")
	defer os.Remove("TestRetained_test.db")
	server, _ := New()
	server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1236",
			DBFilename: "TestRetained_test.db",
		},
	)

	// ###############################################  Write an retained value
	clientIDUUID, _ := uuid.NewRandom()
	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1236")
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })

	// connect
	pahoClient := paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	// subscribe to check that we are connected
	var connected sync.Mutex
	connected.Lock()
	if token := pahoClient.Subscribe("$METRIC/broker/version", 0, func(client paho.Client, msg paho.Message) {
		connected.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	connected.Lock()

	randomUUID := uuid.New()
	if token := pahoClient.Publish("retained/value", 0, true, randomUUID.String()); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClient.Disconnect(200)

	// ###############################################  Check for stored value
	clientIDUUID, _ = uuid.NewRandom()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClient = paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	var subscribeLock sync.Mutex
	subscribeLock.Lock()
	if token := pahoClient.Subscribe("retained/#", 0, func(client paho.Client, msg paho.Message) {
		if string(msg.Payload()) != randomUUID.String() {
			t.FailNow()
			return
		}
		subscribeLock.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	subscribeLock.Lock()

	pahoClient.Disconnect(200)
	server.ServeClose()
}
