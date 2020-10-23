package broker

func (broker *Broker) handleConackPacket(event *Event) (err error) {
	event.client.Connected = true
	return
}
