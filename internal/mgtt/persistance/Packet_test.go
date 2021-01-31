package persistance

import (
	"os"
	"testing"
)

func TestStorePacketIDs(t *testing.T) {
	// delete db if exist
	_, err := os.Stat("TestStorePacketIDs.bolt")
	if err == nil {
		os.Remove("TestStorePacketIDs.bolt")
	}
	Open("TestStorePacketIDs.bolt")
	defer os.Remove("TestStorePacketIDs.bolt")

	packetInfo := PacketInfo{
		Topic:          "Integrationtest",
		OriginClientID: "original",
	}

	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	if packetStoreLastID != packetInfo.MessageID {
		t.FailNow()
	}

	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	if packetStoreLastID != packetInfo.MessageID {
		t.FailNow()
	}

	packetStoreLastID = 10
	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	if packetStoreLastID != 10 {
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

	packetInfo := PacketInfo{
		OriginClientID:  "original",
		OriginMessageID: 5,

		Topic: "Integrationtest",
	}

	if err := PacketStore("TestPacketExist", &packetInfo); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// this should not exist
	var clientIDNotExist = "test"
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &clientIDNotExist,
	}); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this also not
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &clientIDNotExist,
		Topic:          &packetInfo.Topic,
	}); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this also not
	var notExistingOriginMessageID uint16 = 40
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &notExistingOriginMessageID,
		Topic:           &packetInfo.Topic,
	}); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &packetInfo.OriginMessageID,
		Topic:           &packetInfo.Topic,
	}); found == false {
		t.Error(err)
		t.FailNow()
	}
	// this should
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &packetInfo.OriginClientID,
	}); found == false {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _, _ := PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &packetInfo.OriginMessageID,
	}); found == false {
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

	packetInfo := PacketInfo{
		OriginClientID:  "original",
		OriginMessageID: 5,
	}

	// store the first packet
	firstTopic := "First Topic"
	packetInfo.Topic = firstTopic
	if err := PacketStore("TestPacketDelete", &packetInfo); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// store the second packet
	secondTopic := "Second Topic"
	packetInfo.Topic = secondTopic
	if err := PacketStore("TestPacketDelete", &packetInfo); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// try to delete an not existing packet
	var nonExistingClient string = "dummy"
	if err := PacketDelete("TestPacketDelete", PacketFindOpts{
		OriginClientID: &nonExistingClient,
	}); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// delete the first packet
	if err := PacketDelete("TestPacketDelete", PacketFindOpts{
		Topic: &firstTopic,
	}); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// this should
	if found, _, err := PacketExist("TestPacketDelete", PacketFindOpts{
		Topic: &firstTopic,
	}); found == true {
		t.Error(err)
		t.FailNow()
	}

	// this should
	packetInfo.Topic = "Integrationtest"
	if found, _, err := PacketExist("TestPacketDelete", PacketFindOpts{
		Topic: &secondTopic,
	}); found == false {
		t.Error(err)
		t.FailNow()
	}

}
