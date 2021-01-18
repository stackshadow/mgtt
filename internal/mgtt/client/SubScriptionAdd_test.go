package client

import (
	"os"
	"testing"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mocked"
)

func TestSubScriptionAdd(t *testing.T) {
	// setup logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var testClient *MgttClient = &MgttClient{}
	var netserver mocked.Con = mocked.ConNew()
	testClient.Init(netserver, 0)

	testClient.SubScriptionAdd("/test/sensors/#")
	testClient.SubScriptionsAdd([]string{"/test/users/#", "/flat/sensors/#"})

	// send it
	// construct the package
	pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.MessageID = 0
	pub.Retain = false
	pub.TopicName = "/test/sensors/temperature"
	pub.Payload = []byte("")
	pub.Qos = 0
	testClient.Publish(pub)

	readedPacket, err := testClient.PacketRead()
	if err != nil {
		t.FailNow()
	}
	switch respPacket := readedPacket.(type) {
	case *packets.PublishPacket:
		if respPacket.TopicName != "/test/sensors/temperature" {
			t.FailNow()
		}
	default:
		t.FailNow()
	}

	pub = packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.MessageID = 0
	pub.Retain = false
	pub.TopicName = "/test/public/temperature"
	pub.Payload = []byte("")
	pub.Qos = 0
	testClient.Publish(pub)

	pub = packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.MessageID = 0
	pub.Retain = false
	pub.TopicName = "/test/users/admin/write"
	pub.Payload = []byte("")
	pub.Qos = 0
	testClient.Publish(pub)

	readedPacket, err = testClient.PacketRead()
	if err != nil {
		t.FailNow()
	}
	switch respPacket := readedPacket.(type) {
	case *packets.PublishPacket:
		if respPacket.TopicName != "/test/users/admin/write" {
			t.FailNow()
		}
	default:
		t.FailNow()
	}

}
