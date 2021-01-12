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
)

func TestOnSelfUsernameGet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// ############################################### setup auth-password ###############################################

	os.Remove("./TestOnSelfUsernameGet_auth.yml")
	LocalInit("TestOnSelfUsernameGet_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	netserver := connTester{
		packetSendLoopExit: make(chan byte),
	}
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnSelfUsernameGet")
	testClient.UsernameSet("admin")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/password/set/success")
	testClient.SubScriptionAdd("$SYS/self/username/string")
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
			if respPacket.TopicName == "$SYS/auth/user/admin/password/set/success" {
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
			if respPacket.TopicName == "$SYS/self/username/string" &&
				string(respPacket.Payload) == string(respPacket.Payload) {
				respondLock.Unlock()
			} else {
				t.FailNow()
			}
		default:
			t.FailNow()
		}
	}()

	requestLock.Unlock()
	OnHandleMessage("TestOnSelfUsernameGet", "$SYS/auth/user/admin/password/set", []byte("passwordadmintest"))
	respondLock.Lock()

	requestLock.Unlock()
	OnHandleMessage("TestOnSelfUsernameGet", "$SYS/self/username/get", []byte(""))
	respondLock.Lock()

	os.Remove("./TestOnSelfUsernameGet_auth.yml")
}
