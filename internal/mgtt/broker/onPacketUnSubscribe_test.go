package broker

import (
	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test07_UnSubscribe() {

	// connect
	createClient()
	pahoClient, _ := createClient()

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("client7/data", 0, func(client paho.Client, msg paho.Message) {
		suite.FailNow("this should not happen")
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	if token := pahoClient.Unsubscribe("client7/data"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	pahoClient, _ = createClient()
	if token := pahoClient.Publish("client7/data", 0, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect
}
