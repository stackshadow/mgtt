package auth

var adminTopicsEnabled bool

// OnHandleMessage gets called after OnPublishRequest
//
// If this function return true, the plugin handled the message and no other plugin will get it
//
// If a plugin handle the message, it will NOT sended to subscribers
func OnHandleMessage(originClientID string, topic string, payload []byte) (handled bool) {

	// who is currently logged in
	if topic == "$SYS/self/user/get" {
		handled = true
		go publishOwnUserName(originClientID)
	}

	return
}
