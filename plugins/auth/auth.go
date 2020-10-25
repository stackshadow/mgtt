package auth

import "gitlab.com/mgtt/plugin"

func LocalInit() {

	OnInit()

	newPlugin := plugin.V1{
		OnAcceptNewClient: OnAcceptNewClient,
	}
	plugin.Register("auth", &newPlugin)
}

func OnInit() {
	loadConfig("auth.yml")
	go watchConfig()
}

func OnAcceptNewClient(clientID string, username string, password string) (accepted bool) {
	return passwordCheck(username, password)
}
