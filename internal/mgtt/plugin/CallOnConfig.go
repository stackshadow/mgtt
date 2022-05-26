package plugin

import "github.com/rs/zerolog/log"

// CallOnConfig will call the OnConfig-Function on all plugins
func CallOnPluginConfig(pluginConfigs map[string]interface{}) (configChanged bool) {

	for pluginName, plugin := range pluginList {

		pluginData := pluginConfigs[pluginName]
		if pluginData != nil {
			log.Debug().Str("plugin", pluginName).Msg("OnPluginConfig Call")
			configChanged = configChanged || plugin.OnPluginConfig(&pluginData)
			pluginConfigs[pluginName] = pluginData
			log.Debug().Str("plugin", pluginName).Msg("OnPluginConfig Ret")
		} else {
			log.Info().Str("plugin", pluginName).Msg("no config for this plugin present")
		}
	}

	return
}
