package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
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
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1234",
			DBFilename: "TestLastWill_test.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ###############################################  The client with will-message
	clientIDUUID, _ := uuid.NewRandom()
	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1234")
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })

	pahoClientOpts.WillEnabled = true
	pahoClientOpts.WillTopic = "lastwill/value"
	pahoClientOpts.WillRetained = true
	pahoClientOpts.WillQos = 0
	pahoClientOpts.WillPayload = []byte("value")

	// connect and send an retained value
	pahoClient := paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	// ###############################################  The client
	clientIDUUID, _ = uuid.NewRandom()
	pahoClientOpts.SetClientID(clientIDUUID.String())

	// connect and send an retained value
	pahoClientSub := paho.NewClient(pahoClientOpts)
	if token := pahoClientSub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

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
	pahoClient.Disconnect(200)
	lastWillRecvd.Lock()

	// disconnect the last client
	pahoClientSub.Disconnect(200)

	// close the server
	server.ServeClose()
}
