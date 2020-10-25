package plugin

import "github.com/rs/zerolog/log"

// CallOnPublishSendRequest will call the OnPublishRequest-Function on all plugins
func CallOnPublishSendRequest(clientID string, username string, publishTopic string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnPublishSendRequest != nil {
			log.Debug().Str("name", pluginName).Msg("call OnPublishRequest")
			accepted = accepted && plugin.OnPublishSendRequest(clientID, username, publishTopic)
		}
	}

	return
}
