package auth

import (
	"gitlab.com/mgtt/cli"
	"gitlab.com/mgtt/plugin"
)

// LocalInit will init the auth-plugin and register it
func LocalInit() {

	// OnInit open the config file and watch for changes
	OnInit()

	newPlugin := plugin.V1{
		OnAcceptNewClient: OnAcceptNewClient,
	}
	plugin.Register(cli.CLI.ConfigPath+"auth", &newPlugin)
}

// OnInit open the config file and watch for changes
func OnInit() {
	loadConfig(cli.CLI.ConfigPath + "auth.yml")
	go watchConfig()
}

func OnAcceptNewClient(clientID string, username string, password string) (accepted bool) {
	return passwordCheck(username, password)
}
