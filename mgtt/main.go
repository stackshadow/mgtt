package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/cli"
)

func init() {
	// ########################## Command line parse ##########################
	ctx := kong.Parse(&cli.CLI,
		kong.Name("mgtt"),
		kong.Description("Message Go Telemetry Transport"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.Vars{
			"version": cli.Version,
		})
	cli.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped

	// print version
	log.Info().Str("version", cli.Version).Send()

	err := ctx.Run()
	ctx.FatalIfErrorf(err)

}

func main() {
	return
}
