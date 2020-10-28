package plugin

import "github.com/rs/zerolog/log"

// CallOnSubscribeRequest get called when the broker try to publish a message to an subscriber
func CallOnSubscribeRequest(clientID string, username string, publishTopic string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnSubscribeRequest != nil {
			log.Debug().Str("plugin", pluginName).Msg("call OnSubscribeRequest")
			accepted = accepted && plugin.OnSubscribeRequest(clientID, username, publishTopic)
		}
	}

	return
}
