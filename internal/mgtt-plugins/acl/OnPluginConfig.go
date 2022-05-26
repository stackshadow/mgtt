package acl

import (
	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

func OnPluginConfig(pluginData *interface{}) (configChanged bool) {

	// load from interface
	pluginConfigBytes, err := yaml.Marshal(pluginData)
	utils.PanicOnErr(err)
	configLoad(pluginConfigBytes)

	// is the plugin enabled ?
	if !pluginConfig.Enable {
		plugin.DeRegister("acl")
	}

	return
}
