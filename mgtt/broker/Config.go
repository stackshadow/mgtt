package broker

// Config represents the config of your broker
type Config struct {
	Version string

	// the MQTT-Port
	URL string

	TLS      bool
	CAFile   string
	CertFile string
	KeyFile  string

	DBFilename string
}
