package plugin

import "github.com/rs/zerolog/log"

// V1 represents an plugin in version 1
type V1 struct {

	// OnAcceptNewClient gets called, when a new client connects but is not yet added to the list of known clients
	// if this function return false, the client will not added to the known-client-list and get disconnected with return code "not authorized"
	OnAcceptNewClient func(clientID string, username string, password string) bool

	// OnSubscriptionRequest gets called, when an client request an subscription
	// return false if plugin request an abort of an subscription
	OnSubscriptionRequest func(clientID string, username string, subscriptionTopic string) bool

	// OnPublishRequest gets called before an message will be sended to a subscriber
	// return false if plugin request an abort of an publish-message
	OnPublishRequest func(clientID string, username string, publishTopic string) bool

	// OnIncoming get called when an published message arrived from a broker
	OnIncoming func(clientID string, topic string, payload string)
}

var pluginList map[string]*V1 = make(map[string]*V1)

// Register will register a new Plugin
func Register(name string, newPlugin *V1) {
	log.Info().Str("name", name).Msg("Registered new plugin")
	pluginList[name] = newPlugin
}
