package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestSubUnSub(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestSubUnSub_test.db")
	defer os.Remove("TestSubUnSub_test.db")
	server, _ := New()
	  server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1238",
			DBFilename: "TestSubUnSub_test.db",
		},
	)
 

	// vars
	var subscriptionReceived1 sync.Mutex
	var subscriptionReceived2 bool
	subscriptionReceived1.Lock()

	// connect
	clientIDUUID, _ := uuid.NewRandom()
	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1238")
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })

	// subscription-client
	pahoClientSub := paho.NewClient(pahoClientOpts)
	if token := pahoClientSub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	if token := pahoClientSub.Subscribe("home/room1/+", 0, func(client paho.Client, msg paho.Message) {
		subscriptionReceived1.Unlock()
		subscriptionReceived2 = true
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	// publishing-client
	clientIDUUID, _ = uuid.NewRandom()
	pahoClientOpts.SetClientID(clientIDUUID.String())

	pahoClientPub := paho.NewClient(pahoClientOpts)
	if token := pahoClientPub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	// publish
	if token := pahoClientPub.Publish("home/room1/light1", 1, true, "on"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	subscriptionReceived1.Lock()

	// unsubscribe
	subscriptionReceived2 = false
	if token := pahoClientPub.Unsubscribe("home/room1/+"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)
	if subscriptionReceived2 == true {
		t.Fail()
	}

	pahoClientPub.Disconnect(500)
	pahoClientSub.Disconnect(500)

	server.ServeClose()
	time.Sleep(time.Second * 1)
}
