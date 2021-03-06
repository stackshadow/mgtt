package broker

// Current is the last created broker
var Current *Broker = nil

// New will create a new Broker
func New() (broker *Broker, err error) {
	broker = &Broker{
		clientEvents:                make(chan *Event, 10),
		pubrecs:                     make(map[uint16]Qos2),
		loopHandleResendPacketsExit: make(chan bool),
	}

	// remember the current broker
	Current = broker

	return
}
