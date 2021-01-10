package plugin

import "github.com/rs/zerolog/log"

// CallOnDisconnected get called, when an client is disconnected
func CallOnDisconnected(clientID string) {

	for pluginName, plugin := range pluginList {
		if plugin.OnDisconnected != nil {
			log.Debug().Str("plugin", pluginName).Msg("call OnDisconnected")
			plugin.OnDisconnected(clientID)
		}
	}

	return
}
