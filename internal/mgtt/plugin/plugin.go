package plugin

import (
	"github.com/rs/zerolog/log"
)

// V1 represents an plugin in version 1
type V1 struct {

	// OnConfig gets called when the app loads or the configfile was changed
	OnPluginConfig func(*interface{}) (configChanged bool)

	// OnNewClient gets called when a new client is incoming
	OnNewClient func(remoteAddr string)

	// OnAcceptNewClient gets called, when a CONNECT-Packet arrived but is not yet added to the list of known clients
	//
	// if this function return false, the client will not added to the known-client-list and get disconnected with return code "not authorized"
	OnAcceptNewClient func(clientID string, username string, password string) bool

	// OnConnected will called after connection is established from an subscriber on the broker
	OnConnected func(clientID string)

	// OnDisconnected called, when a client leaves our list of clients
	OnDisconnected func(clientID string)

	// OnConnack will called after a conack was received
	OnConnack func(clientID string)

	// OnSubscriptionRequest gets called, when an client request an subscription
	// return false if plugin request an abort of an subscription
	OnSubscriptionRequest func(clientID string, username string, subscriptionTopic string) bool

	// OnPublishRequest get called when an publisher try to publish to the broker
	//
	// If an plugin return false, connection will be closed
	OnPublishRequest func(clientID string, username string, topic string) bool

	// OnHandleMessage gets called after OnPublishRequest
	//
	// If this function return true, the plugin handled the message and no other plugin will get it
	//
	// If a plugin handle the message, it will NOT sended to subscribers
	OnHandleMessage func(originClientID string, topic string, payload []byte) (handled bool)

	// OnSendToSubscriberRequest get called when the broker try to publish a message to an subscriber
	//
	// if an plugin set it to false, the message will NOT be sended to clientID
	//
	// This function gets called BEFORE check if the subscriber subscribe to the topic
	//
	// clientID: The clientID of the target client
	// username: The username of the target client
	// publishTopic: The topic the broker try to publish to the subscriber
	OnSendToSubscriberRequest func(clientID string, username string, publishTopic string) bool
}

var pluginList map[string]*V1 = make(map[string]*V1)

// Register will register a new Plugin
func Register(name string, newPlugin *V1) {
	log.Info().Str("name", name).Msg("Registered new plugin")
	pluginList[name] = newPlugin
}

// DeRegister will remove an plugin from the pluginlist
func DeRegister(name string) {
	log.Info().Str("name", name).Msg("DeRegister plugin")
	delete(pluginList, name)
}
