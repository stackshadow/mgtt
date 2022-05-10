package broker

import (
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteTests(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// setup suite
func (suite *TestSuite) SetupSuite() {

	var err error
	config.Values.Level = "debug"
	config.Values.URL = "tcp://127.0.0.1:1235"
	config.Values.DB = "./integrationtest.db"
	config.Apply()

	// Broker
	var newbroker *Broker
	newbroker, err = New()
	utils.PanicOnErr(err)

	go newbroker.Serve()

	time.Sleep(time.Second * 2)
}

func (suite *TestSuite) BeforeTest(suiteName, testName string) {
	log.Debug().Msg("#################### START " + testName + " ####################")
}

func (suite *TestSuite) AfterTest(suiteName, testName string) {
	log.Debug().Msg("#################### STOP " + testName + " ####################\n")

}

func (suite *TestSuite) TearDownSuite() {

}

func createClient(clientID ...string) (pahoClient paho.Client, clientUUID string) {

	if len(clientID) > 0 {
		clientUUID = clientID[0]
	} else {
		clientUUID = uuid.Must(uuid.NewRandom()).String()
	}

	var pahoClientConnected sync.Mutex
	pahoClientConnected.Lock()

	pahoClientOpts := paho.NewClientOptions()
	pahoClientOpts.SetClientID(clientUUID)
	pahoClientOpts.SetUsername("dummy")
	pahoClientOpts.SetPassword("dummy")
	pahoClientOpts.AddBroker(config.Values.URL)
	pahoClientOpts.SetAutoReconnect(true)
	pahoClientOpts.SetOnConnectHandler(func(c paho.Client) { pahoClientConnected.Unlock() })
	pahoClientOpts.SetWill("lastwill/client", pahoClientOpts.ClientID, 0, false)
	pahoClientOpts.SetCleanSession(false)

	// connect
	pahoClient = paho.NewClient(pahoClientOpts)
	if token := pahoClient.Connect(); token.Wait() && token.Error() != nil {
		utils.PanicOnErr(token.Error())
	}
	pahoClientConnected.Lock() // wait for connected

	// subscribe to check that we are connected
	var connected sync.Mutex
	connected.Lock()
	if token := pahoClient.Subscribe("$METRIC/broker/version", 0, func(client paho.Client, msg paho.Message) {
		connected.Unlock()
	}); token.Wait() && token.Error() != nil {
		utils.PanicOnErr(token.Error())
	}
	connected.Lock()

	//
	if token := pahoClient.Subscribe("clients/ping", 0, func(client paho.Client, msg paho.Message) {
		if token := client.Publish("clients/pong", 0, false, pahoClientOpts.ClientID); token.Wait() && token.Error() != nil {
			utils.PanicOnErr(token.Error())
		}
	}); token.Wait() && token.Error() != nil {
		utils.PanicOnErr(token.Error())
	}

	return pahoClient, pahoClientOpts.ClientID
}
