package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
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
	serverConfig := Config{
		URL:        "tcp://127.0.0.1:1235",
		DBFilename: "TestPubSubQoS1_test.db",
	}
	  server.Serve(serverConfig)
 

	// ###############################################  The client with will-message
	var pahoClientConnected bool = false

	pahoClientSub, err := testConnectClient(serverConfig.URL, "", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// vars
	var subscriptionReceivedQoS1 sync.Mutex

	if token := pahoClientSub.Subscribe("qos/1", 0, func(client paho.Client, msg paho.Message) {
		subscriptionReceivedQoS1.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	// publishing-client
	pahoClientPub, err := testConnectClient(serverConfig.URL, "", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// publish QoS1
	subscriptionReceivedQoS1.Lock()
	if token := pahoClientPub.Publish("qos/1", 1, true, "100%"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)
	subscriptionReceivedQoS1.Lock()

	pahoClientPub.Disconnect(500)
	pahoClientSub.Disconnect(500)

	server.ServeClose()
	time.Sleep(time.Second * 1)
}
