package cli

// CLICommon holds common stuff
type CLICommon struct {
	// Debug will enable debug mode
	Debug    DebugFlag    `help:"Enable debug mode." short:"v" env:"DEBUG" default:"false"`
	Terminal TerminalFlag `help:"Enable terminal mode ( log not printed as json)" env:"TERMINAL" default:"false"`
}

// CLI is the overall cli-struct
var CLI struct {
	CLICommon

	ConfigPath string `help:"Path where config files are stored. This can be used by plugins"  env:"CONFIGPATH" default:"./"`

	URL      string `help:"Set the url where the broker is listening"  env:"URL" default:"tcp://0.0.0.0:8883"`
	CertFile string `help:"The certificate to use for TLS. If not set, TLS is disabled"  env:"CERT" default:"cert.pem"`
	KeyFile  string `help:"The private key to use for TLS"  env:"KEY" default:"key.pem"`

	DBFilename string `help:"Filename for retained message-db"  env:"DBFILENAME" default:"messages.db"`

	// ConnectTimeout holds the timeout in seconds for CONNECT
	ConnectTimeout int64 `help:"Timeout in seconds for CONNECT. If an client don't send a connect after this time, it will be disconnected" env:"CONNECT_TIMEOUT" default:"30"`
}
