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

func TestOnAuthUserPasswordSet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Remove("./TestOnAuthUserPasswordSet_auth.yml")
	LocalInit("TestOnAuthUserPasswordSet_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	var netserver = connTester{
		packetSendLoopExit: make(chan byte, 1024),
	}
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnAuthUserPasswordSet")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/password/set/success")
	clientlist.Add(testClient)

	var respondLock sync.Mutex
	respondLock.Lock()

	go func() {
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
	}()

	OnHandleMessage("TestOnAuthUserPasswordSet", "$SYS/auth/user/admin/password/set", []byte("admin"))
	respondLock.Lock()

	// check correct password
	if OnAcceptNewClient("TestOnAuthUserPasswordSet", "admin", "admin") != true {
		t.Fail()
	}

	// check wrong password
	if OnAcceptNewClient("TestOnAuthUserPasswordSet", "admin", "admin2") == true {
		t.Fail()
	}

	os.Remove("./TestOnAuthUserPasswordSet_auth.yml")
}
