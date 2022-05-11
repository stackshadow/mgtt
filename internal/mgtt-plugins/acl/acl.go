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
func Init() {

	// init the map
	config.Plugins.ACL.Rules = make(map[string][]aclEntry)

	// create the plugin
	newPlugin := plugin.V1{
		OnConfig:                  OnConfig,
		OnPublishRequest:          OnPublishRequest,
		OnSendToSubscriberRequest: OnSendToSubscriberRequest,
	}

	// register it
	plugin.Register("acl", &newPlugin)
}
