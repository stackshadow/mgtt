package plugin

import "github.com/rs/zerolog/log"

// V1 represents an plugin in version 1
type V1 struct {

	// OnAcceptNewClient gets called, when a new client connects but is not yet added to the list of known clients
	// if this function return false, the client will not added to the known-client-list and get disconnected with return code "not authorized"
	OnAcceptNewClient func(clientID string, username string, password string) bool

	// OnConnected will called after connection is established from an subscriber on the broker
	OnConnected func(clientID string)

	// OnConnack will called after a conack was received
	OnConnack func(clientID string)

	// OnSubscriptionRequest gets called, when an client request an subscription
	// return false if plugin request an abort of an subscription
	OnSubscriptionRequest func(clientID string, username string, subscriptionTopic string) bool

	// OnPublishRecvRequest get called when an published message arrived on the broker from a publisher
	// return false to abort publish to all subscribers
	OnPublishRecvRequest func(clientID string, topic string, payload string) bool

	// OnPublishSendRequest gets called before an message will be sended to a subscriber
	// return false to abort publish to an specific subscriber
	OnPublishSendRequest func(clientID string, username string, publishTopic string) bool
}

var pluginList map[string]*V1 = make(map[string]*V1)

// Register will register a new Plugin
func Register(name string, newPlugin *V1) {
	log.Info().Str("name", name).Msg("Registered new plugin")
	pluginList[name] = newPlugin
}
