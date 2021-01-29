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

func TestPubSubQoS1(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestPubSubQoS1_test.db")
	defer os.Remove("TestPubSubQoS1_test.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1235",
			DBFilename: "TestPubSubQoS1_test.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ###############################################  The client with will-message
	clientIDUUID := uuid.New()
	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()
	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientIDUUID.String())
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker("tcp://127.0.0.1:1235")
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })

	// vars
	var subscriptionReceivedQoS1 sync.Mutex
	var subscriptionReceivedQoS2 sync.Mutex
	subscriptionReceivedQoS1.Lock()
	subscriptionReceivedQoS2.Lock()

	// subscription-client
	pahoClientSub := paho.NewClient(pahoClientOpts)
	if token := pahoClientSub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClientConnected.Lock() // wait for connected

	if token := pahoClientSub.Subscribe("qos/1", 0, func(client paho.Client, msg paho.Message) {
		subscriptionReceivedQoS1.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	if token := pahoClientSub.Subscribe("qos/2", 0, func(client paho.Client, msg paho.Message) {
		subscriptionReceivedQoS2.Unlock()
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

	// publish QoS1

	if token := pahoClientPub.Publish("qos/1", 1, true, "100%"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)
	subscriptionReceivedQoS1.Lock()

	if token := pahoClientPub.Publish("qos/2", 2, true, "100%"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)
	subscriptionReceivedQoS2.Lock()

	pahoClientPub.Disconnect(500)
	pahoClientSub.Disconnect(500)

	server.ServeClose()
	time.Sleep(time.Second * 3)
}
