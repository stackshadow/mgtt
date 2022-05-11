package plugin

// CallOnConnack will call the OnConnack-Function on all plugins
func CallOnConnack(clientID string) {

	for _, plugin := range pluginList {
		if plugin.OnConnack != nil {
			plugin.OnConnack(clientID)
		}
	}

}
