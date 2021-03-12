package broker

import (
	"os"
	"testing"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestPubSys(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestPubSys.db")
	defer os.Remove("TestPubSys.db")
	server, _ := New()
	serverConfig := Config{
		URL:        "tcp://127.0.0.1:1239",
		DBFilename: "TestPubSys.db",
	}
	server.Serve(serverConfig)

	// ###############################################  a client that subscribe
	var pahoClientConnected bool = false

	pahoClient, err := testConnectClient(serverConfig.URL, "SysSubscribeClient", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// we subscribe to several topics
	if token := pahoClient.Subscribe("$SYS/#", 0, func(client paho.Client, msg paho.Message) {
		t.Errorf("Another client received $SYS.. that should not happen")
		t.Fail()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	// ###############################################  a client that publish
	pahoClientPublish, err := testConnectClient(serverConfig.URL, "SysPublishClient", false, &pahoClientConnected)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// we subscribe to several topics
	if token := pahoClientPublish.Publish("$SYS/getPassword", 0, false, ""); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	pahoClient.Disconnect(200)
	pahoClientPublish.Disconnect(200)

	server.ServeClose()
}
