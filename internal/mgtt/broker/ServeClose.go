package broker

import "gitlab.com/mgtt/internal/mgtt/clientlist"

func (b *Broker) ServeClose() {

	b.loopHandleResendPacketsExit <- true
	clientlist.RemoveAll()
	b.serverListener.Close()

}
