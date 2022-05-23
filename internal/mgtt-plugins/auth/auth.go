package auth

import (
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

// LocalInit will init the auth-plugin and register it
func Init() {

	// init the map
	pluginConfig.Users = make(map[string]pluginConfigUser)

	// create the plugin
	newPlugin := plugin.V1{
		OnPluginConfig:    OnPluginConfig,
		OnAcceptNewClient: OnAcceptNewClient,
		OnHandleMessage:   OnHandleMessage,
	}

	// register it
	plugin.Register("auth", &newPlugin)
}

// OnAcceptNewClient gets called, when a CONNECT-Packet arrived but is not yet added to the list of known clients
func OnAcceptNewClient(clientID string, username string, password string) (accepted bool) {
	return configCheckPassword(username, password)
}
