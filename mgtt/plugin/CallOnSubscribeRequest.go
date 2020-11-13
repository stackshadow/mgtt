package plugin

import "github.com/rs/zerolog/log"

// CallOnSendToSubscriberRequest get called when the broker try to publish a message to an subscriber
func CallOnSendToSubscriberRequest(clientID string, username string, publishTopic string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnSendToSubscriberRequest != nil {
			log.Debug().Str("plugin", pluginName).Msg("call OnSendToSubscriberRequest")
			accepted = accepted && plugin.OnSendToSubscriberRequest(clientID, username, publishTopic)
		}
	}

	return
}
