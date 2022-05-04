package persistance

func (suite *TestSuite) TestPacketDelete() {

	var err error

	packetInfo := PacketInfo{
		OriginClientID:  "original",
		OriginMessageID: 5,
	}

	// store the first packet
	firstTopic := "First Topic"
	packetInfo.Topic = firstTopic
	err = PacketStore("TestPacketDelete", &packetInfo)
	suite.NoError(err)

	// store the second packet
	secondTopic := "Second Topic"
	packetInfo.Topic = secondTopic
	err = PacketStore("TestPacketDelete", &packetInfo)
	suite.NoError(err)

	// try to delete an not existing packet
	var nonExistingClient string = "dummy"
	err = PacketDelete("TestPacketDelete", PacketFindOpts{
		OriginClientID: &nonExistingClient,
	})
	suite.NoError(err)

	// delete the first packet
	err = PacketDelete("TestPacketDelete", PacketFindOpts{
		Topic: &firstTopic,
	})
	suite.NoError(err)

	// this should
	var found bool

	// this should not exist
	found, _, err = PacketExist("TestPacketDelete", PacketFindOpts{
		Topic: &firstTopic,
	})
	suite.False(found)
	suite.NoError(err)

	// this should exist
	packetInfo.Topic = "Integrationtest"
	found, _, err = PacketExist("TestPacketDelete", PacketFindOpts{
		Topic: &secondTopic,
	})
	suite.True(found)
	suite.NoError(err)
}
