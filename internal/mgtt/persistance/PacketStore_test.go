package persistance

func (suite *TestSuite) TestStore() {

	packetInfo := PacketInfo{
		Topic:          "Integrationtest",
		OriginClientID: "original",
	}

	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	suite.Assert().Equal(packetStoreLastID, packetInfo.MessageID)

	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	suite.Assert().Equal(packetStoreLastID, packetInfo.MessageID)

	packetStoreLastID = 10
	packetInfo.MessageID = 0
	PacketStore("TestStorePacketIDs", &packetInfo)
	suite.Assert().Equal(packetStoreLastID, uint16(10))

	var counter = 0
	PacketIterate("TestStorePacketIDs", func(info PacketInfo) {
		counter++
	})

	suite.Assert().Equal(3, counter)
}
