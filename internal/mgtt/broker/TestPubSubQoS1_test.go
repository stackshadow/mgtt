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

func ServeForTests(t *testing.T) {

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
			URL:        "tcp://127.0.0.1:1237",
			DBFilename: "test1.db",
		},
	)
	time.Sleep(time.Millisecond * 500)
}

func TestPubSubQoS1(t *testing.T) {

	// We serve the server
	ServeForTests(t)

	// vars
	var subscriptionReceivedQoS1 sync.Mutex
	var subscriptionReceivedQoS2 sync.Mutex
	subscriptionReceivedQoS1.Lock()
	subscriptionReceivedQoS2.Lock()

	// connect
	clientIDUUID, _ := uuid.NewRandom()
	pahoClientSubOpts := paho.NewClientOptions()
	pahoClientSubOpts.SetClientID(clientIDUUID.String())
	pahoClientSubOpts.SetUsername("dummy")
	pahoClientSubOpts.SetPassword("dummy")
	pahoClientSubOpts.AddBroker("tcp://127.0.0.1:1237")
	pahoClientSubOpts.SetAutoReconnect(true)

	// subscription-client
	pahoClientSub := paho.NewClient(pahoClientSubOpts)
	if token := pahoClientSub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
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
	pahoClientPubOpts := paho.NewClientOptions()
	pahoClientPubOpts.SetClientID(clientIDUUID.String())
	pahoClientPubOpts.SetUsername("dummy")
	pahoClientPubOpts.SetPassword("dummy")
	pahoClientPubOpts.AddBroker("tcp://127.0.0.1:1237")
	pahoClientPubOpts.SetAutoReconnect(true)

	pahoClientPub := paho.NewClient(pahoClientPubOpts)
	if token := pahoClientPub.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

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
}
