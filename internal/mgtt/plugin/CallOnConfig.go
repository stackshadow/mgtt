package plugin

// CallOnConfig will call the OnConfig-Function on all plugins
func CallOnConfig(yamlConfigData []byte) {

	for _, plugin := range pluginList {
		if plugin.OnConfig != nil {
			plugin.OnConfig(yamlConfigData)
		}
	}

}
