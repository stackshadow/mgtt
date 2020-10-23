package plugin

import "github.com/rs/zerolog/log"

// CallOnPublishRequest will call the OnPublishRequest-Function on all plugins
//
// Only if all plugins return true, the publish is accepted. This can mainly used for ACL-Plugins
func CallOnPublishRequest(clientID string, username string, publishTopic string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnPublishRequest != nil {
			log.Debug().Str("name", pluginName).Msg("call OnPublishRequest")
			accepted = accepted && plugin.OnPublishRequest(clientID, username, publishTopic)
		}
	}

	return
}
