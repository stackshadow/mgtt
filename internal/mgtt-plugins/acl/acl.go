package acl

import "gitlab.com/mgtt/internal/mgtt/plugin"

type aclEntry struct {
	Route string `yaml:"route"`

	// "r" / "w"
	Direction string `yaml:"direction"`

	// true if allow, false if not
	Allow bool `yaml:"allow"`
}

// LocalInit will init the auth-plugin and register it
func LocalInit(ConfigPath string) {

	// OnInit open the config file and watch for changes
	OnInit(ConfigPath)

	newPlugin := plugin.V1{
		OnPublishRequest:          OnPublishRequest,
		OnSendToSubscriberRequest: OnSendToSubscriberRequest,
	}
	plugin.Register("acl", &newPlugin)
}

// OnInit open the config file and watch for changes
func OnInit(ConfigPath string) {
	loadConfig(ConfigPath + "acl.yml")
	go watchConfig()
}
