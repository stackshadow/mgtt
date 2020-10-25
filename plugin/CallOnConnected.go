package plugin

// CallOnConnected will call the OnPublishRecvRequest-Function on all plugins
func CallOnConnected(clientID string) {

	for _, plugin := range pluginList {
		if plugin.OnConnected != nil {
			plugin.OnConnected(clientID)
		}
	}

	return
}
