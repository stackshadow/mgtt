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

func TestOnAuthUserSet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// ############################################### setup auth-password ###############################################

	os.Remove("./TestOnSelfUsernameGet_auth.yml")
	LocalInit("TestOnSelfUsernameGet_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	var netserver mocked.Con = mocked.ConNew()
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnSelfUsernameGet")
	testClient.UsernameSet("admin")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/set/success")
	testClient.SubScriptionAdd("$SYS/self/user/json")
	testClient.SubScriptionAdd("error")
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
			if respPacket.TopicName == "error" {
				respondLock.Unlock()
				return
			}
			t.FailNow()
		}

	}()

	requestLock.Unlock()
	OnHandleMessage("TestOnSelfUsernameGet", "$SYS/auth/user//set", []byte("{ \"password\": \"admin\" }"))
	respondLock.Lock()

	os.Remove("./TestOnSelfUsernameGet_auth.yml")
}
