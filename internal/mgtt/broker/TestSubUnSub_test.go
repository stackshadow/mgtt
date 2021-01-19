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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// We serve the server
	ServeForTests(t)

	// vars
	var subscriptionReceived1 sync.Mutex
	var subscriptionReceived2 bool
	subscriptionReceived1.Lock()

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
	if token := pahoClientSub.Subscribe("home/room1/+", 0, func(client paho.Client, msg paho.Message) {
		subscriptionReceived1.Unlock()
		subscriptionReceived2 = true
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

	testserver.ServeClose()
	time.Sleep(time.Second * 1)
}
