package messagestore

import (
	"os"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func Test_Store(t *testing.T) {

	// delete db if exist
	_, err := os.Stat("integrtest.db")
	if err == nil {
		os.Remove("integrtest.db")
	}

	store, err := Open("integrtest.db")
	if err != nil {
		t.Error(err)
	}

	// store an packet
	newControlPacket := packets.NewControlPacket(packets.Publish)
	newPacket := newControlPacket.(*packets.PublishPacket)
	newPacket.TopicName = "Integrationtest"

	var newMessageID uint16 = 0
	option := StoreResendPacketOptions{
		ResendAt: time.Now().Add(time.Minute * 1),
		Packet:   newPacket,
	}

	err = store.StoreResendPacket("integrtest", &newMessageID, &option)
	if newMessageID != 0 || err != nil {
		t.FailNow()
	}

	// try again
	newMessageID = 0
	err = store.StoreResendPacket("integrtest", &newMessageID, &option)
	if newMessageID != 1 || err != nil {
		t.FailNow()
	}

	// try again
	newMessageID = 5
	err = store.StoreResendPacket("integrtest", &newMessageID, &option)
	if newMessageID != 5 || err != nil {
		t.FailNow()
	}

	// try again
	newMessageID = 0
	err = store.StoreResendPacket("integrtest", &newMessageID, &option)
	if newMessageID != 2 || err != nil {
		t.FailNow()
	}

	// iterate, there should be two packages
	var counter int
	store.IterateResendPackets("integrtest", func(storedInfo *StoreResendPacketOptions) {
		counter++
		if storedInfo.Packet.TopicName != "Integrationtest" {
			t.FailNow()
		}
	})
	if counter != 4 {
		t.FailNow()
	}

}
