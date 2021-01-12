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

func TestOnAuthUserDelete(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Remove("./TestOnAuthUserDelete_auth.yml")
	LocalInit("TestOnAuthUserDelete_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	netserver := connTester{
		packetSendLoopExit: make(chan byte),
	}
	testClient.Init(netserver, 0)
	testClient.IDSet("integrationtest")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/deleteme/password/set/success")
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
			if respPacket.TopicName == "$SYS/auth/user/deleteme/password/set/success" {
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
	OnHandleMessage("integrationtest", "$SYS/auth/user/deleteme/password/set", []byte("admin"))
	respondLock.Lock()

	requestLock.Unlock()
	OnHandleMessage("integrationtest", "$SYS/auth/user/deleteme/delete", []byte("admin"))
	respondLock.Lock()

	os.Remove("./TestOnAuthUserDelete_auth.yml")
}
