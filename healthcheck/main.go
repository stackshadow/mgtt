package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
)

var cli struct {
	URI string `help:"URI to check" env:"HEALTHURI" default:"http://localhost:8080/health"`
}

func init() {
	// ########################## Command line parse ##########################
	kong.Parse(&cli,
		kong.Name("health"),
		kong.Description("Health check for distroless-container"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
}

func main() {
	if healthy(cli.URI) == false {
		log.Println("NOT OK")
		os.Exit(1)
	}
	log.Println("OK")
}

func healthy(uiURI string) bool {
	if resp, err := http.Get(uiURI); err == nil {
		defer resp.Body.Close()
		return true
	}
	return false
}
