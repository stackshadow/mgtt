package client

import (
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mocked"
)

func TestSend(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.Logger.With().Caller().Logger()

	var testClient *MgttClient = &MgttClient{}
	var netserver mocked.Con = mocked.ConNew()
	testClient.Init(netserver, 0)

	testClient.SubScriptionAdd("/test/sensors/#")
	testClient.SubScriptionsAdd([]string{"/test/users/#", "/flat/sensors/#"})

	// ConnackPacket
	err := testClient.SendConnack(1, false)
	if err != nil {
		t.FailNow()
	}

	readedPacket, err := testClient.PacketRead()
	xType := reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.ConnackPacket" {
		t.FailNow()
	}

	// ConnectPacket
	testClient.SendConnect("username", "password", "clientid")

	readedPacket, err = testClient.PacketRead()
	xType = reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.ConnectPacket" {
		t.FailNow()
	}

	// PingreqPacket
	testClient.SendPingreq()

	readedPacket, err = testClient.PacketRead()
	xType = reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.PingreqPacket" {
		t.FailNow()
	}

	// PingrespPacket
	testClient.SendPingresp()

	readedPacket, err = testClient.PacketRead()
	xType = reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.PingrespPacket" {
		t.FailNow()
	}

	// PubackPacket
	err = testClient.SendPuback(1)
	if err != nil {
		t.FailNow()
	}

	readedPacket, err = testClient.PacketRead()
	xType = reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.PubackPacket" {
		t.FailNow()
	}

	// PubcompPacket
	err = testClient.SendPubcomp(1)
	if err != nil {
		t.FailNow()
	}

	readedPacket, err = testClient.PacketRead()
	xType = reflect.TypeOf(readedPacket)
	if xType.String() != "*packets.PubcompPacket" {
		t.FailNow()
	}
}
