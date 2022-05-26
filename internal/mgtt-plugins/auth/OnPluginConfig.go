package auth

import (
	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

func OnPluginConfig(pluginData *interface{}) (configChanged bool) {

	// load from interface
	pluginConfigBytes, err := yaml.Marshal(pluginData)
	utils.PanicOnErr(err)
	configChanged = configLoad(pluginConfigBytes)

	// is the plugin enabled ?
	if !pluginConfig.Enable {
		plugin.DeRegister("auth")
		return
	}

	// we need to store it back
	if configChanged {
		var tempData []byte
		tempData, err = yaml.Marshal(pluginConfig)
		utils.PanicOnErr(err)
		yaml.Unmarshal(tempData, pluginData)
	}

	return
}
