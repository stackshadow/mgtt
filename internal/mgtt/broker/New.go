package broker

// New will create a new Broker
func New() (broker *Broker, err error) {
	broker = &Broker{
		clientEvents:                make(chan *Event, 10),
		pubrecs:                     make(map[uint16]Qos2),
		loopHandleResendPacketsExit: make(chan bool),
	}

	return
}
