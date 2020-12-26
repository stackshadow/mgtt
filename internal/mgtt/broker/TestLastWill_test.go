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
	os.Remove("test1.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1238",
			DBFilename: "test1.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ###############################################  The client with will-message
	clientIDUUID := uuid.New()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1238")
	pahoClientOpts.SetAutoReconnect(true)

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

	// ###############################################  The client
	pahoClientSubUUID := uuid.New()
	pahoClientSubOpts := paho.NewClientOptions()
	pahoClientSubOpts.SetClientID(pahoClientSubUUID.String())
	pahoClientSubOpts.SetUsername("dummy")
	pahoClientSubOpts.SetPassword("dummy")
	pahoClientSubOpts.AddBroker("tcp://127.0.0.1:1238")
	pahoClientSubOpts.SetAutoReconnect(true)

	// connect and send an retained value
	pahoClientSub := paho.NewClient(pahoClientSubOpts)
	if token := pahoClientSub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
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
	pahoClient.Disconnect(200)
	lastWillRecvd.Lock()

	// disconnect the last client
	pahoClientSub.Disconnect(200)

	// close the server
	server.ServeClose()
	time.Sleep(time.Second * 3)
}
