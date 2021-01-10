package auth

import (
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
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

func TestOnHandleUserGet(t *testing.T) {
	os.Remove("./integrationtest_auth.yml")
	LocalInit("integrationtest_")

	// create a dummy client
	var testClient *client.MgttClient = &client.MgttClient{}
	netserver := connTester{
		packetSendLoopExit: make(chan byte, 10),
	}
	testClient.Init(netserver, 0)
	testClient.IDSet("integrationtest")
	testClient.Connected = true
	testClient.SubScriptionAdd("$SYS/auth/user/admin/error")
	testClient.SubScriptionAdd("$SYS/auth/user/admin/json")
	testClient.SubScriptionAdd("$SYS/auth/user/admin/password/set/success")
	clientlist.Add(testClient)

	var respondLock sync.Mutex
	respondLock.Lock()

	go func() {
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
	}()

	// wait for the error
	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/get", []byte(""))
	time.Sleep(time.Millisecond * 100)
	respondLock.Lock()

	// add the user, it should now exist
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

	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/password/set", []byte("admin"))
	time.Sleep(time.Millisecond * 100)
	respondLock.Lock()

	// add the user, it should now exist
	go func() {
		respondPacket, _ := testClient.PacketRead()
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

	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/get", []byte(""))
	time.Sleep(time.Millisecond * 100)
	respondLock.Lock()

	os.Remove("./integrationtest_auth.yml")

}
