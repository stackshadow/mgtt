package auth

import (
	"gitlab.com/mgtt/plugin"
)

// LocalInit will init the auth-plugin and register it
func LocalInit(ConfigPath string) {

	// OnInit open the config file and watch for changes
	OnInit(ConfigPath)

	newPlugin := plugin.V1{
		OnAcceptNewClient: OnAcceptNewClient,
	}
	plugin.Register("auth", &newPlugin)
}

// OnInit open the config file and watch for changes
func OnInit(ConfigPath string) {
	loadConfig(ConfigPath + "auth.yml")
	go watchConfig()
}

// OnAcceptNewClient gets called, when a CONNECT-Packet arrived but is not yet added to the list of known clients
func OnAcceptNewClient(clientID string, username string, password string) (accepted bool) {
	return passwordCheck(username, password)
}
