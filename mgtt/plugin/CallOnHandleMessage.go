package plugin

// CallOnHandleMessage call all OnHandleMessage-functions
//
// If this function return true, the next plugin will handle the message
//
// If a plugin handle the message, it will NOT sended to subscribers
func CallOnHandleMessage(originClientID string, topic string, payload []byte) (messageWasHandled bool) {

	for _, plugin := range pluginList {
		if plugin.OnHandleMessage != nil {

			if plugin.OnHandleMessage(originClientID, topic, payload) == false {
				break
			} else {
				messageWasHandled = true
			}

		}
	}

	return
}
