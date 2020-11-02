package main

import (
	"github.com/alecthomas/kong"
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
			"version": "0.0.1",
		})
	cli.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped
	err := ctx.Run()
	ctx.FatalIfErrorf(err)

}

func main() {
	return
}
