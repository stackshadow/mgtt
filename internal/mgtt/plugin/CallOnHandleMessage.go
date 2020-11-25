package plugin

// CallOnHandleMessage call all OnHandleMessage-functions
//
// If this function return true, the plugin handled the message and no other plugin will get it
//
// If a plugin handle the message, it will NOT sended to subscribers
func CallOnHandleMessage(originClientID string, topic string, payload []byte) (messageWasHandled bool) {

	for _, plugin := range pluginList {
		if plugin.OnHandleMessage != nil {

			if plugin.OnHandleMessage(originClientID, topic, payload) == true {
				messageWasHandled = true
				break
			}

		}
	}

	return
}
