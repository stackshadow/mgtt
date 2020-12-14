package messagestore

import (
	"os"
	"testing"
	"time"
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

	option := PacketInfo{
		ResendAt: time.Now().Add(time.Minute * 1),
		Topic:    "Integrationtest",
	}

	err = store.StoreResendPacket("integrtest", &option)
	if option.MessageID != 0 || err != nil {
		t.Error("newMessageID != 0")
		t.FailNow()
	}

	// try again
	option.MessageID = 0
	err = store.StoreResendPacket("integrtest", &option)
	if option.MessageID != 1 || err != nil {
		t.Error("newMessageID != 1")
		t.FailNow()
	}

	// try again
	option.MessageID = 5
	err = store.StoreResendPacket("integrtest", &option)
	if option.MessageID != 5 || err != nil {
		t.Error("newMessageID != 5")
		t.FailNow()
	}

	// try again
	option.MessageID = 0
	err = store.StoreResendPacket("integrtest", &option)
	if option.MessageID != 2 || err != nil {
		t.Error("newMessageID != 2")
		t.FailNow()
	}

	// iterate, there should be two packages
	var counter int
	store.IterateResendPackets("integrtest", func(storedInfo *PacketInfo) {
		counter++
		if storedInfo.Topic != "Integrationtest" {
			t.Error("storedInfo.Topic != 'Integrationtest'")
			t.FailNow()
		}
	})
	if counter != 4 {
		t.Error("counter != 4")
		t.FailNow()
	}

}
