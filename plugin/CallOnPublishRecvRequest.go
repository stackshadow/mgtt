package plugin

import "github.com/rs/zerolog/log"

// CallOnPublishRecvRequest will call the OnPublishRecvRequest-Function on all plugins
func CallOnPublishRecvRequest(clientID string, topic string, payload string) (accepted bool) {

	//per default we should accept
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnPublishRecvRequest != nil {
			log.Debug().Str("name", pluginName).Msg("call OnIncoming")
			accepted = accepted && plugin.OnPublishRecvRequest(clientID, topic, payload)
		}
	}

	return
}
