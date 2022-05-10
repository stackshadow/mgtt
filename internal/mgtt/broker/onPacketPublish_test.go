package broker

import (
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test02_Publish() {

	// connect
	pahoClient, _ := createClient()

	var published sync.Mutex
	published.Lock()
	publishedCounter := 0

	// we subscribe as client 1 and respond with the data we get
	if token := pahoClient.Subscribe("client1/echo/respond", 0, func(client paho.Client, msg paho.Message) {
		if string(msg.Payload()) == "test2" {
			published.Unlock()
			publishedCounter++
			return
		}

	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	if token := pahoClient.Publish("client1/echo/write2", 0, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect
	if token := pahoClient.Publish("client1/echo/write", 0, false, "test2"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	published.Lock()

	suite.Assert().Equal(1, publishedCounter)

	// we should have 2 clients
	// because the broker not send a message back to us
	clientCounter := 0
	if token := pahoClient.Subscribe("clients/pong", 0, func(client paho.Client, msg paho.Message) {
		clientCounter++
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}
	if token := pahoClient.Publish("clients/ping", 0, false, ""); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	time.Sleep(time.Second * 2)
	suite.Assert().Equal(2, clientCounter)

}

func (suite *TestSuite) Test03_PublishRetained() {

	// connect
	pahoClient, _ := createClient()

	var published sync.Mutex
	published.Lock()

	if token := pahoClient.Publish("retained/data/one", 0, true, "retained"); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	} // this should not have any effect

	pahoClient, _ = createClient()
	if token := pahoClient.Subscribe("retained/#", 0, func(client paho.Client, msg paho.Message) {
		suite.Assert().Equal("retained/data/one", string(msg.Topic()))
		published.Unlock()
		return
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	published.Lock()
}
