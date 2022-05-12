package broker

import (
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test05_QoS2() {

	// connect
	pahoClient, _ := createClient()

	var published sync.Mutex
	published.Lock()

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("qos2/data", 2, func(client paho.Client, msg paho.Message) {
		suite.Assert().Equal("test2", string(msg.Payload()))
		published.Unlock()
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	pahoClient, _ = createClient()
	if token := pahoClient.Publish("qos2/data", 2, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect

	published.Lock()
}
