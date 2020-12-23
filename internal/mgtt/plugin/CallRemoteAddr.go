package plugin

// CallOnNewClient will call the OnNewClient-Function on all plugins
func CallOnNewClient(remoteAddr string) {

	for _, plugin := range pluginList {
		if plugin.OnNewClient != nil {
			plugin.OnNewClient(remoteAddr)
		}
	}

	return
}
