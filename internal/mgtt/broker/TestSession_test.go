package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func testConnectClient(url string, clientID string, cleanSession bool, connected *bool) (pahoClient paho.Client, err error) {

	if clientID == "" {
		clientIDUUID, _ := uuid.NewRandom()
		clientID = clientIDUUID.String()
	}

	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientID)
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker(url)
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) {
		*connected = true
	})
	pahoClientOpts.SetCleanSession(cleanSession)

	pahoClientOpts.WillEnabled = true
	pahoClientOpts.WillTopic = "lastwill/value"
	pahoClientOpts.WillRetained = true
	pahoClientOpts.WillQos = 0
	pahoClientOpts.WillPayload = []byte("value")

	// connect
	*connected = false
	pahoClient = paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}
	// wait for connected
	for {
		if *connected == true {
			break
		}
		time.Sleep(time.Millisecond)
	}

	/*
		// subscribe to check that we are connected
		var onBrokerVersionLock sync.Mutex
		onBrokerVersionLock.Lock()

		if token := pahoClient.Subscribe("$SYS/broker/version", 0, func(client paho.Client, msg paho.Message) {
			onBrokerVersionLock.Unlock()
		}); token.Wait() && token.Error() != nil {
			err = token.Error()
			return
		}
		onBrokerVersionLock.Lock()
	*/
	return
}

func TestSession(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()
	/*
		log.Logger = log.Logger.
			Output(zerolog.ConsoleWriter{Out: &lumberjack.Logger{
				Filename:   "./foo.log",
				MaxSize:    1,  // megabytes after which new file is created
				MaxBackups: 3,  // number of backups
				MaxAge:     28, //days
			}}).
			With().
			Caller().
			Logger()
	*/

	// ############################################### the broker
	os.Remove("TestSession_test.db")
	defer os.Remove("TestSession_test.db")
	server, _ := New()
	serverConfig := Config{
		URL:        "tcp://127.0.0.1:1237",
		DBFilename: "TestSession_test.db",
	}
	go server.Serve(serverConfig)
	time.Sleep(time.Second * 1)

	// ###############################################  Write an retained value
	var pahoClientConnected bool = false

	pahoClient, err := testConnectClient(serverConfig.URL, "sessionClient", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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
	time.Sleep(time.Second)

	// ###############################################  Connect again
	pahoClient, err = testConnectClient(serverConfig.URL, "sessionClient", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var published sync.Mutex
	pahoClient.AddRoute("sensors/room2/#", func(client paho.Client, msg paho.Message) {
		published.Unlock()
	})
	published.Lock()

	// ###############################################  Connect a publisher
	pahoClientPublish, err := testConnectClient(serverConfig.URL, "publisherClient", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if token := pahoClientPublish.Publish("sensors/room2/#", 0, true, "100%"); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	published.Lock()

	pahoClient.Disconnect(200)
	time.Sleep(time.Second)

	server.ServeClose()
}
