package persistance

func (suite *TestSuite) TestPacketExist() {

	var err error

	packetInfo := PacketInfo{
		OriginClientID:  "original",
		OriginMessageID: 5,

		Topic: "Integrationtest",
	}

	err = PacketStore("TestPacketExist", &packetInfo)
	suite.NoError(err)

	// this should not exist
	var found bool
	var clientIDNotExist = "test"
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &clientIDNotExist,
	})
	suite.False(found)
	suite.NoError(err)

	// this also not
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &clientIDNotExist,
		Topic:          &packetInfo.Topic,
	})
	suite.False(found)
	suite.NoError(err)

	// this also not
	var notExistingOriginMessageID uint16 = 40
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &notExistingOriginMessageID,
		Topic:           &packetInfo.Topic,
	})
	suite.False(found)
	suite.NoError(err)

	// this should
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &packetInfo.OriginMessageID,
		Topic:           &packetInfo.Topic,
	})
	suite.True(found)
	suite.NoError(err)

	// this should
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginClientID: &packetInfo.OriginClientID,
	})
	suite.True(found)
	suite.NoError(err)

	// this should
	found, _, err = PacketExist("TestPacketExist", PacketFindOpts{
		OriginMessageID: &packetInfo.OriginMessageID,
	})
	suite.True(found)
	suite.NoError(err)
}
