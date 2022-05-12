package cli

import (
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/mgtt/internal/mgtt/server"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

// Common holds common stuff
type Common struct {
	// Debug will enable debug mode
	Config string `help:"Enable debug mode." short:"c"   default:"./mgtt.conf"`
}

// CLI is the overall cli-struct
var cliData Common
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
			"version": config.Version,
		},
	)

	// load the config
	config.MustLoad(cliData.Config)

	// print version
	log.Info().
		Str("version", config.Version).
		Msg("Welcome to mgtt")
}

// Run will execute Commands
func Run() {

	var err error

	// TLS
	if config.Values.TLS.CA.File != "" {
		server.MustCreateCA()
	}

	// Cert
	if config.Values.TLS.Cert.File != "" {
		server.MustCreateCert()
	}

	// Broker
	var newbroker *broker.Broker
	newbroker, err = broker.New()
	utils.PanicOnErr(err)

	var done chan bool
	done, err = newbroker.Serve()
	if err != nil {
		log.Error().Err(err).Send()
	}

	<-done

}

func isTest() bool {

	for _, flag := range os.Args {
		if strings.Contains(flag, "-test.") {
			return true
		}

	}
	return false
}
