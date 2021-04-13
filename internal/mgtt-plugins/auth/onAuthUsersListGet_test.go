package auth

import (
	"os"
	"sync"
	"testing"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mocked"
)

func TestOnAuthUsersListGet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Setenv("ENABLE_ADMIN_TOPICS", "true")
	os.Remove("./TestOnAuthUsersListGet_auth.yml")
	LocalInit("TestOnAuthUsersListGet_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	var netserver mocked.Con = mocked.ConNew()
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnAuthUsersListGet")
	testClient.UsernameSet("admin")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/set/success")
	testClient.SubScriptionAdd("$SYS/auth/users/list/json")
	clientlist.Add(testClient)

	var requestLock sync.Mutex
	var respondLock sync.Mutex
	requestLock.Lock()
	respondLock.Lock()

	go func() {
		requestLock.Lock()
		respondPacket, _ := testClient.PacketRead()
		switch respPacket := respondPacket.(type) {
		case *packets.PublishPacket:
			if respPacket.TopicName == "$SYS/auth/users/list/json" {
				respondLock.Unlock()
			} else {
				t.FailNow()
			}
		default:
			t.FailNow()
		}

	}()

	// delete user
	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUsersListGet", "$SYS/auth/users/list/get", []byte(""))
	respondLock.Lock()

	os.Remove("./TestOnAuthUsersListGet_auth.yml")
}
