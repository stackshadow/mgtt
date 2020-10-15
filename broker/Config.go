package broker

// Config represents the config of your broker
type Config struct {
	// the MQTT-Port
	URL string

	// the MQTTS-Port
	URLTLS string
}
