package plugin

import "github.com/rs/zerolog/log"

// CallOnPublishRequest get called when an publisher try to publish to the broker
func CallOnPublishRequest(clientID string, username string, topic string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnPublishRequest != nil {
			log.Debug().Str("plugin", pluginName).Msg("call OnPublishRequest")
			accepted = accepted && plugin.OnPublishRequest(clientID, username, topic)
		}
	}

	return
}
