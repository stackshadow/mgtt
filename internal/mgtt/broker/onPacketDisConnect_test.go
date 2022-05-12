package broker

import (
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test06_LastWill() {

	// connect
	pahoClientThatDisconnect, disconnectedClientID := createClient()
	pahoClient, _ := createClient()

	var published sync.Mutex
	published.Lock()

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("lastwill/client", 0, func(client paho.Client, msg paho.Message) {
		suite.Assert().Equal(disconnectedClientID, string(msg.Payload()))
		published.Unlock()
		return
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	if token := pahoClientThatDisconnect.Publish("client4/data", 1, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect
	pahoClientThatDisconnect.Disconnect(500)

	published.Lock()
	pahoClient.Disconnect(5000)
}
