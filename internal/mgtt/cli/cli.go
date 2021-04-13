package cli

import (
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

// Common holds common stuff
type Common struct {
	// Debug will enable debug mode
	Debug    DebugFlag    `help:"Enable debug mode." short:"v" env:"DEBUG" default:"false"`
	Terminal TerminalFlag `help:"Enable terminal mode ( log not printed as json)" env:"TERMINAL" default:"false"`
	CreateEnvHelpFileCommand
}

// Parameter represents the parameter of your CLI
type Parameter struct {
	Common

	CreateCA   CmdCreateCA   `kong:"cmd,help='Create a ca'"`
	CreateCert CmdCreateCert `kong:"cmd,help:'Create a cert'"`
	Serve      CmdServe      `kong:"cmd,help:'Serve'"`

	ConfigPath string `help:"Path where config files are stored. This can be used by plugins"  env:"CONFIGPATH" default:"./"`

	// ConnectTimeout holds the timeout in seconds for CONNECT
	ConnectTimeout int64 `help:"Timeout in seconds for CONNECT. If an client don't send a connect after this time, it will be disconnected" env:"CONNECT_TIMEOUT" default:"30"`

	RetryInMinutes int64 `help:"Retry delay of failed QoS1 QoS2 in Minutes" env:"RETRY" default:"1"`

	Plugins string `help:"Name of enabled plugins comma separated" env:"PLUGINS" default:"auth,acl"`

	EnableAdminTopics AdminTopicsFlag `kong:"cmd,help='Enable admin topics',env='ENABLE_ADMIN_TOPICS',default='false'"`
}

// CLI is the overall cli-struct
var cliData Parameter
var ctxKong *kong.Context

// ParseAndRun will parse command line arguments and run commands
func init() {

	// we not init on cmd-line-test
	if isTest() == true {
		return
	}

	// ########################## Command line parse ##########################
	ctxKong = kong.Parse(&cliData, // this trigger AfterApply()
		kong.Name("mgtt"),
		kong.Description("Message Go Telemetry Transport"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.Vars{
			"version": Version,
		},
	)

	ctxKong.ApplyDefaults()
	ctxKong.Bind(cliData)

	// print version
	log.Info().Str("version", Version).Send()
}

// Run will execute Commands
func Run() {
	err := ctxKong.Run()
	ctxKong.FatalIfErrorf(err)
}

func isTest() bool {
	for _, flag := range os.Args {
		if strings.Contains(flag, "-test.run") || strings.Contains(flag, "-test.v") {
			return true
		}

	}
	return false
}
