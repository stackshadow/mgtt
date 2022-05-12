package broker

import (
	"sync"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func (suite *TestSuite) Test08_Session() {

	// connect
	pahoClient, pahoClientID := createClient()

	var published sync.Mutex
	published.Lock()

	if token := pahoClient.Subscribe("client/sessions/thisnot", 0, func(client paho.Client, msg paho.Message) {
		suite.FailNow("this should not happen")
		return
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	if token := pahoClient.Subscribe("client/sessions/data", 0, func(client paho.Client, msg paho.Message) {
		published.Unlock()
		return
	}); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	// we disconnect this client to use the session afterwards
	pahoClient.Disconnect(5000)

	// this should not work
	pahoPublishClient, _ := createClient()
	if token := pahoPublishClient.Publish("client/sessions/thisnot", 0, false, ""); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	// we reconnect with our clientID to get the session back
	pahoClient, pahoClientID = createClient(pahoClientID)
	pahoClient.AddRoute("client/sessions/data", func(client paho.Client, msg paho.Message) {
		published.Unlock()
		return
	})

	// tro to send it
	if token := pahoPublishClient.Publish("client/sessions/data", 0, false, ""); token.Wait() && token.Error() != nil {
		suite.NoError(token.Error())
	}

	published.Lock()
}
