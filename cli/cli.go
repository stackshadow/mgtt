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

	URL      string `help:"Set the url where the broker is listening"  env:"URL" default:"tcp://0.0.0.0:1883"`
	CertFile string `help:"The certificate to use for TLS. If not set, TLS is disabled"  env:"CERT" default:"cert.pem"`
	KeyFile  string `help:"The private key to use for TLS"  env:"KEY" default:"key.pem"`
}
