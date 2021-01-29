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
	log.Logger = log.Logger.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Caller().
		Logger()

	// ############################################### the broker
	os.Remove("TestRetained_test.db")
	defer os.Remove("TestRetained_test.db")
	server, _ := New()
	go server.Serve(
		Config{
			URL:        "tcp://127.0.0.1:1238",
			DBFilename: "TestRetained_test.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ###############################################  Write an retained value
	pahoOpts, _ := TESTCreatePahoOpts()
	somerandomvalue, _ := uuid.NewRandom()

	// connect
	pahoClient := paho.NewClient(pahoOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

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

	if token := pahoClient.Publish("retained/value", 0, true, somerandomvalue.String()); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	pahoClient.Disconnect(200)

	// ###############################################  Check for stored value

	// connect again and we should get the value
	var subscribeLock sync.Mutex
	pahoOpts, _ = TESTCreatePahoOpts()
	pahoClient = paho.NewClient(pahoOpts)

	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

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

	pahoClient.Disconnect(200)
	server.ServeClose()
}
