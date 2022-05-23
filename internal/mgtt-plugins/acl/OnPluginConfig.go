package acl

import (
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func OnPluginConfig(yamlConfigData []byte) (configChanged bool) {
	configLoad(yamlConfigData)

	// is the plugin enabled ?
	if !pluginConfig.Enable {
		plugin.DeRegister("acl")
	}

	return
}
