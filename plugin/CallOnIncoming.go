package plugin

import "github.com/rs/zerolog/log"

// CallOnIncoming will call the OnIncoming-Function on all plugins
func CallOnIncoming(clientID string, topic string, payload string) {

	for pluginName, plugin := range pluginList {
		if plugin.OnPublishRequest != nil {
			log.Debug().Str("name", pluginName).Msg("call OnIncoming")
			plugin.OnIncoming(clientID, topic, payload)
		}
	}

	return
}
