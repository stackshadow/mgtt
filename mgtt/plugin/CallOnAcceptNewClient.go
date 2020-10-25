package plugin

import "github.com/rs/zerolog/log"

// CallOnAcceptNewClient will call the OnAcceptNewClient-Function on all plugins
//
// Only if all plugins return true, the connection is accepted
func CallOnAcceptNewClient(clientID string, username string, password string) (accepted bool) {

	//per default we should connect
	accepted = true

	for pluginName, plugin := range pluginList {
		if plugin.OnAcceptNewClient != nil {
			log.Debug().Str("name", pluginName).Msg("call OnAcceptNewClient")
			accepted = accepted && plugin.OnAcceptNewClient(clientID, username, password)
		}
	}

	return
}
