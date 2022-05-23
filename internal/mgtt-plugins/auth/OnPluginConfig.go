package auth

import (
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

func OnPluginConfig(yamlConfigData []byte) (configChanged bool) {
	configChanged = configLoad(yamlConfigData)
	if configChanged {
		config.AlterPluginConfig("auth", pluginConfig)
	}

	// is the plugin enabled ?
	if !pluginConfig.Enable {
		plugin.DeRegister("auth")
	}

	return
}
