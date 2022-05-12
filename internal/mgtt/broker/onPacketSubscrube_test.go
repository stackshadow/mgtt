package broker

import (
	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test01_Subscribe() {

	// connect
	createClient()
	pahoClient, _ := createClient()

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("client1/echo/write", 0, func(client paho.Client, msg paho.Message) {
		if token := client.Publish("client1/echo/respond", 0, false, msg.Payload()); token.Wait() && token.Error() != nil {
			suite.NoError(token.Error())
		}
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}
}
