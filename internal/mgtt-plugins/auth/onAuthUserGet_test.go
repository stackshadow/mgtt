package auth

import (
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

type connTester struct {
	packetSendLoopExit chan byte
}

func (c connTester) Read(b []byte) (n int, err error) {
	bytes := <-c.packetSendLoopExit
	b[0] = bytes
	return 1, nil
}
func (c connTester) Write(b []byte) (n int, err error) {
	for _, singleByte := range b {
		c.packetSendLoopExit <- singleByte
	}

	return 1, nil
}

func (c connTester) Close() error {
	return nil
}

func (c connTester) LocalAddr() (add net.Addr) {
	return
}

func (c connTester) RemoteAddr() (add net.Addr) {
	return
}

func (c connTester) SetDeadline(t time.Time) error {
	return nil
}

func (c connTester) SetReadDeadline(t time.Time) error {
	return nil
}

func (c connTester) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestOnAuthUserGet(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	os.Remove("./TestOnAuthUserGet_auth.yml")
	LocalInit("TestOnAuthUserGet_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	netserver := connTester{
		packetSendLoopExit: make(chan byte),
	}
	testClient.Init(netserver, 0)
	testClient.IDSet("TestOnAuthUserGet")
	testClient.UsernameSet("admin")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/error")
	testClient.SubScriptionAdd("$SYS/auth/user/admin/json")
	testClient.SubScriptionAdd("$SYS/auth/user/admin/set/success")
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
			if respPacket.TopicName == "$SYS/auth/user/admin/error" {
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
			if respPacket.TopicName == "$SYS/auth/user/admin/set/success" {
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
			if respPacket.TopicName == "$SYS/auth/user/admin/json" {
				respondLock.Unlock()
			} else {
				t.FailNow()
			}
		default:
			t.FailNow()
		}

	}()

	// wait for the error
	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUserGet", "$SYS/auth/user/admin/get", []byte(""))
	respondLock.Lock()

	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUserGet", "$SYS/auth/user/admin/set", []byte("{ \"password\": \"admin\" }"))
	respondLock.Lock()

	requestLock.Unlock()
	OnHandleMessage("TestOnAuthUserGet", "$SYS/auth/user/admin/get", []byte(""))
	respondLock.Lock()

	os.Remove("./TestOnAuthUserGet_auth.yml")
}
