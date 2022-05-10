package broker

import (
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test04_Publish() {

	// connect
	pahoClient, _ := createClient()

	var published sync.Mutex
	published.Lock()

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("client4/data", 0, func(client paho.Client, msg paho.Message) {
		suite.Assert().Equal("test2", string(msg.Payload()))
		published.Unlock()
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	pahoClient, _ = createClient()
	if token := pahoClient.Publish("client4/data", 1, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect

	published.Lock()
}
