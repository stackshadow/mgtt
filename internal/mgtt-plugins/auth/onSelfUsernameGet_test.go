package auth

import (
	"os"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
)

func TestOnSelfUsernameGet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// ############################################### setup auth-password ###############################################

	os.Remove("./integrationtest_auth.yml")
	LocalInit("integrationtest_")

	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/password/set", []byte("admin"))
	time.Sleep(time.Millisecond * 500)

	// create a plugin that unlocks the connection of a client
	var clientLock sync.Mutex

	// ############################################### the broker ###############################################

	os.Remove("test1.db")
	server, _ := broker.New()
	go server.Serve(
		broker.Config{
			URL:        "tcp://127.0.0.1:1212",
			DBFilename: "test1.db",
		},
	)
	time.Sleep(time.Second * 1)

	// ############################################### client ###############################################
	pahoOpts := paho.NewClientOptions()

	clientIDUUID, _ := uuid.NewRandom()
	pahoOpts.SetClientID(clientIDUUID.String())
	pahoOpts.SetUsername("admin")
	pahoOpts.SetPassword("admin")
	pahoOpts.AddBroker("tcp://127.0.0.1:1212")
	pahoOpts.SetAutoReconnect(true)

	pahoClient := paho.NewClient(pahoOpts)

	clientLock.Lock()
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}

	if token := pahoClient.Subscribe("$SYS/self/username", 0, func(client paho.Client, msg paho.Message) {
		if string(msg.Payload()) == "admin" {
			clientLock.Unlock()
		} else {
			t.Error("Payload not contain username")
		}

	}); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)

	if token := pahoClient.Publish("$SYS/self/username/get", 0, true, ""); token.Wait() && token.Error() != nil {
		t.Error(token.Error())
		t.FailNow()
	}
	time.Sleep(time.Second * 1)
	pahoClient.Disconnect(200)

	clientLock.Lock()
}
