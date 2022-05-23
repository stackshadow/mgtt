package plugin

import "gopkg.in/yaml.v3"

// CallOnConfig will call the OnConfig-Function on all plugins
func CallOnPluginConfig(pluginConfig map[string]interface{}) (configChanged bool) {

	for pluginName, plugin := range pluginList {
		if plugin.OnPluginConfig != nil {

			out, err := yaml.Marshal(pluginConfig[pluginName])
			if err == nil {
				configChanged = configChanged || plugin.OnPluginConfig(out)
			}

		}
	}

	return
}
