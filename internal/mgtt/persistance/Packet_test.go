package persistance

import (
	"os"
	"testing"
	"time"
)

func TestStorePacketIDs(t *testing.T) {
	// delete db if exist
	_, err := os.Stat("TestStorePacketIDs.bolt")
	if err == nil {
		os.Remove("TestStorePacketIDs.bolt")
	}
	Open("TestStorePacketIDs.bolt")
	defer os.Remove("TestStorePacketIDs.bolt")

	var lastID uint16
	option := PacketInfo{
		ResendAt:       time.Now().Add(time.Minute * 1),
		MessageID:      100,
		Topic:          "Integrationtest",
		PubComp:        false,
		OriginClientID: "original",
	}

	PacketStore(option, &lastID)
	if lastID != 0 {
		t.FailNow()
	}

	lastID = 0
	PacketStore(option, &lastID)
	if lastID != 1 {
		t.FailNow()
	}

	lastID = 10
	PacketStore(option, &lastID)
	if lastID != 10 {
		t.FailNow()
	}
}

func TestPacketExist(t *testing.T) {
	// delete db if exist
	_, err := os.Stat("TestStorePacketExist.bolt")
	if err == nil {
		os.Remove("TestStorePacketExist.bolt")
	}
	Open("TestStorePacketExist.bolt")
	defer os.Remove("TestStorePacketExist.bolt")

	var lastID uint16 = 1
	option := PacketInfo{
		ResendAt:       time.Now().Add(time.Minute * 1),
		MessageID:      100,
		Topic:          "Integrationtest",
		PubComp:        false,
		OriginClientID: "original",
	}

	if err := PacketStore(option, &lastID); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// this should not exist
	var clientIDNotExist = "test"
	if found, _ := PacketExist(&clientIDNotExist, nil, nil); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this also not
	if found, _ := PacketExist(&clientIDNotExist, &option.Topic, nil); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this also not
	if found, _ := PacketExist(&clientIDNotExist, nil, &lastID); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _ := PacketExist(nil, &option.Topic, &lastID); found == false {
		t.Error(err)
		t.FailNow()
	}
	// this should
	if found, _ := PacketExist(&option.OriginClientID, nil, nil); found == false {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _ := PacketExist(nil, nil, &lastID); found == false {
		t.Error(err)
		t.FailNow()
	}

}

func TestPacketDelete(t *testing.T) {
	// delete db if exist
	_, err := os.Stat("TestPacketDelete.bolt")
	if err == nil {
		os.Remove("TestPacketDelete.bolt")
	}
	Open("TestPacketDelete.bolt")
	defer os.Remove("TestPacketDelete.bolt")

	var lastID uint16 = 1
	option := PacketInfo{
		ResendAt:       time.Now().Add(time.Minute * 1),
		MessageID:      100,
		Topic:          "Integrationtest",
		PubComp:        false,
		OriginClientID: "original",
	}

	if err := PacketStore(option, &lastID); err != nil {
		t.Error(err)
		t.FailNow()
	}

	option.Topic = "Second Topic"
	if err := PacketStore(option, &lastID); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// try to delete an not existing packet
	var nonExistingClient string = "dummy"
	if err := PacketDelete(&nonExistingClient, nil, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _ := PacketExist(nil, nil, &lastID); found == false {
		t.Error(err)
		t.FailNow()
	}

	// delete the first packet
	if err := PacketDelete(&option.OriginClientID, nil, nil); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _ := PacketExist(&option.OriginClientID, nil, nil); found == false {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _ := PacketExist(nil, &option.Topic, nil); found == false {
		t.Error(err)
		t.FailNow()
	}
}

func TestPubComp(t *testing.T) {

	// delete db if exist
	_, err := os.Stat("TestPubComp.bolt")
	if err == nil {
		os.Remove("TestPubComp.bolt")
	}
	Open("TestPubComp.bolt")
	defer os.Remove("TestPubComp.bolt")

	var lastID uint16 = 1
	option := PacketInfo{
		ResendAt:       time.Now().Add(time.Minute * 1),
		MessageID:      100,
		Topic:          "Integrationtest",
		PubComp:        false,
		OriginClientID: "original",
	}

	if err := PacketStore(option, &lastID); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if isSet, _, _ := PacketPubCompIsSet(lastID); isSet == true {
		t.Error("PubComp is set, but it should not")
		t.FailNow()
	}

	// set it
	if err := PacketPubCompSet(lastID, true); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if isSet, _, _ := PacketPubCompIsSet(100); isSet == false {
		t.Error("PubComp is not set, but it should")
		t.FailNow()
	}
}
