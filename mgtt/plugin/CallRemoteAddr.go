package plugin

// CallRemoteAddr will call the OnPublishRecvRequest-Function on all plugins
func CallRemoteAddr(remoteAddr string) {

	for _, plugin := range pluginList {
		if plugin.OnNewClient != nil {
			plugin.OnNewClient(remoteAddr)
		}
	}

	return
}
