package broker

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/plugin"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

/*
	cmd := exec.Command(
		"mosquitto_pub",
		"-L",
		"mqtts://admin:admin@127.0.0.1:1234/$SYS/broker/cr",
		"-m",
		"2",
		"-d",
		"-q",
		"0",
	)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
*/

func TESTCreatePahoOpts() (*paho.ClientOptions, error) {
	opts := paho.NewClientOptions()

	clientIDUUID, _ := uuid.NewRandom()
	opts.SetClientID(clientIDUUID.String())
	opts.SetUsername("dummy")
	opts.SetPassword("dummy")
	opts.AddBroker("tcp://127.0.0.1:1238")
	opts.SetAutoReconnect(true)
	return opts, nil
}

func TestRetained(t *testing.T) {

	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// create a plugin that unlocks the connection of a client
	var clientLock sync.Mutex
	var clientLockPlugin plugin.V1
	clientLockPlugin.OnConnected = func(clientID string) {
		clientLock.Unlock()
	}
	plugin.Register("clientLock", &clientLockPlugin)

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

	// ###############################################  Write an retained value
	pahoOpts, _ := TESTCreatePahoOpts()
	somerandomvalue, _ := uuid.NewRandom()

	// connect and send an retained value
	pahoClient := paho.NewClient(pahoOpts)

	clientLock.Lock()
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	clientLock.Lock() // the plugin should unlock this

	if token := pahoClient.Publish("retained/value", 0, true, somerandomvalue.String()); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClient.Disconnect(200)
	time.Sleep(time.Second * 1)

	// ###############################################  Check for stored value

	// connect again and we should get the value
	var subscribeLock sync.Mutex
	pahoClient = paho.NewClient(pahoOpts)

	clientLock.Unlock()
	clientLock.Lock()
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	clientLock.Lock() // the plugin should unlock this

	subscribeLock.Lock()
	if token := pahoClient.Subscribe("retained/#", 0, func(client paho.Client, msg paho.Message) {
		if string(msg.Payload()) != somerandomvalue.String() {
			t.FailNow()
			return
		}
		subscribeLock.Unlock()
	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	subscribeLock.Lock()
	time.Sleep(time.Second * 1)
	pahoClient.Disconnect(200)
}
