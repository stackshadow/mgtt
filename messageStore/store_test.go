package messagestore

import (
	"os"
	"testing"

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
	var newPacket packets.PublishPacket

	newID, err := store.StorePacketWithID("integrtest", 0, &newPacket)
	if newID != 0 || err != nil {
		t.FailNow()
	}

	// try again
	newID, err = store.StorePacketWithID("integrtest", 0, &newPacket)
	if newID != 1 || err != nil {
		t.FailNow()
	}

	// reserve an ID
	var reservedID uint16
	reservedID, err = store.StorePacketWithID("integrtest", 0, nil)
	if reservedID != 2 || err != nil {
		t.FailNow()
	}

	// and we store again a package
	newID, err = store.StorePacketWithID("integrtest", 0, &newPacket)
	if newID != 3 || err != nil {
		t.FailNow()
	}

	// we store the reserved
	newID, err = store.StorePacketWithID("integrtest", reservedID, &newPacket)
	if newID != reservedID || err != nil {
		t.FailNow()
	}

}
