package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/broker"
	"gitlab.com/mgtt/cli"
	parameter "gitlab.com/mgtt/cli"
)

func init() {
	// ########################## Command line parse ##########################
	kong.Parse(&cli.CLI,
		kong.Name("mgtt"),
		kong.Description("Message Go Telemetry Transport"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	parameter.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped
}

func main() {

	newBroker, err := broker.Serve(
		broker.Config{
			URL: ":8883",
		},
	)
	newBroker.Communicate()

	if err != nil {
		log.Error().Err(err).Send()
	}
}
