package broker

import "gitlab.com/mgtt/internal/mgtt/clientlist"

// ServeClose will close all client-connections and broker-listeners
func (b *Broker) ServeClose() {

	// b.loopHandleResendPacketsExit <- true
	clientlist.RemoveAll()
	b.serverListener.Close()

}
