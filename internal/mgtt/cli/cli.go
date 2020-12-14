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

	CreateCA   CmdCreateCA   `cmd help:"Create a ca"`
	CreateCert CmdCreateCert `cmd help:"Create a cert"`
	Serve      CmdServe      `cmd help:"Serve"`

	ConfigPath string `help:"Path where config files are stored. This can be used by plugins"  env:"CONFIGPATH" default:"./"`

	// ConnectTimeout holds the timeout in seconds for CONNECT
	ConnectTimeout int64 `help:"Timeout in seconds for CONNECT. If an client don't send a connect after this time, it will be disconnected" env:"CONNECT_TIMEOUT" default:"30"`

	Plugins string `help:"Name of enabled plugins comma separated"  env:"PLUGINS" default:"auth,acl"`
}