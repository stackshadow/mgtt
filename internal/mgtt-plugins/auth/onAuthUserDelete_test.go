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

func TestOnAuthUserDelete(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Setenv("ENABLE_ADMIN_TOPICS", "true")
	os.Remove("./TestOnAuthUserDelete_auth.yml")
	LocalInit("TestOnAuthUserDelete_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	var netserver mocked.Con = mocked.ConNew()
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnAuthUserDelete")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/deleteme/set/success")
	testClient.SubScriptionAdd("$SYS/auth/user/deleteme/delete/success")
	clientlist.Add(testClient)

	var requestLock sync.Mutex
	var respondLock sync.Mutex
	requestLock.Lock()
	respondLock.Lock()

	// add the user, it should now exist
	go func() {
		requestLock.Lock()
		respondPacket, _ := testClient.PacketRead()
		switch respPacket := respondPacket.(type) {
		case *packets.PublishPacket:
			if respPacket.TopicName == "$SYS/auth/user/deleteme/set/success" {
				respondLock.Unlock()
			} else {
				t.FailNow()
			}
		default:
			t.FailNow()
		}

		requestLock.Lock()
		respondPacket, _ = testClient.PacketRead()
		switch respPacket := respondPacket.(type) {
		case *packets.PublishPacket:
			if respPacket.TopicName == "$SYS/auth/user/deleteme/delete/success" {
				respondLock.Unlock()
			} else {
				t.FailNow()
			}
		default:
			t.FailNow()
		}
	}()

	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUserDelete", "$SYS/auth/user/deleteme/set", []byte("{ \"password\": \"admin\" }"))
	respondLock.Lock()

	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUserDelete", "$SYS/auth/user/deleteme/delete", []byte(""))
	respondLock.Lock()

	os.Remove("./TestOnAuthUserDelete_auth.yml")
}
