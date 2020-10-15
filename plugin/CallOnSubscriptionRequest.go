package plugin

import "github.com/rs/zerolog/log"

// CallOnSubscriptionRequest will call the OnAcceptNewClient-Function on all plugins
//
// Only if all plugins return true, the connection is accepted
func CallOnSubscriptionRequest(clientID string, subscriptions string) (accepted bool) {
	//per default we should connect
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnSubscriptionRequest != nil {
			log.Debug().Str("name", pluginName).Msg("call OnConnect")
			accepted = accepted && plugin.OnSubscriptionRequest(clientID, subscriptions)
		}
	}

	return
}
