package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func TestSession(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestSession_test.db")
	defer os.Remove("TestSession_test.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1237",
			DBFilename: "TestSession_test.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ###############################################  Write an retained value
	clientIDUUID, _ := uuid.NewRandom()
	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1237")
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })
	pahoClientOpts.SetCleanSession(false)

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
	if token := pahoClient.Subscribe("$SYS/broker/version", 0, func(client paho.Client, msg paho.Message) {
		connected.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	connected.Lock()

	// we subscribe to several topics
	if token := pahoClient.Subscribe("sensors/room1/#", 0, func(client paho.Client, msg paho.Message) {

	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	if token := pahoClient.Subscribe("sensors/room2/#", 0, func(client paho.Client, msg paho.Message) {

	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClient.Disconnect(200)

	// ###############################################  Connect again
	// connect
	pahoClient = paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	var published sync.Mutex
	pahoClient.AddRoute("sensors/room2/#", func(client paho.Client, msg paho.Message) {
		published.Unlock()
	})
	published.Lock()

	// ###############################################  Connect a publisher
	pahoClientOpts.SetClientID("clientPublisher")
	pahoClientPublish := paho.NewClient(pahoClientOpts)
	if token := pahoClientPublish.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	if token := pahoClientPublish.Publish("sensors/room2/#", 0, true, "100%"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	published.Lock()

	pahoClient.Disconnect(200)
	server.ServeClose()
}
