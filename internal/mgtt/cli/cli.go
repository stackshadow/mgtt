package cli

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt-plugins/acl"
	"gitlab.com/mgtt/internal/mgtt-plugins/auth"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/mgtt/internal/mgtt/server"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Common holds common stuff
type Common struct {
	// Debug will enable debug mode
	Config string `help:"Enable debug mode." short:"c" default:""`

	Password string `help:"password to bcrypt" short:"p"`
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

	// plugins
	acl.Init()
	auth.Init()

	// config
	configAsJSON, _ := os.LookupEnv("CONFIG_JSON")
	config.MustLoadFromFile(cliData.Config)
	config.MustLoadFromString(configAsJSON)
	configChanged := plugin.CallOnPluginConfig(config.Globals.Plugins) // inform all plugins about the config change
	config.ApplyDefaults()
	config.ApplyLog()
	if configChanged {
		config.MustSave()
	}

	// print version
	log.Info().
		Str("version", config.Version).
		Msg("Welcome to mgtt")
}

// Run will execute Commands
func Run() {

	var err error

	// username / password
	if cliData.Password != "" {
		var bcryptedData []byte
		bcryptedData, err = bcrypt.GenerateFromPassword([]byte(cliData.Password), bcrypt.DefaultCost)
		base64String := base64.StdEncoding.EncodeToString(bcryptedData)
		fmt.Print(base64String)
		os.Exit(0)
	}

	// TLS
	if config.Globals.TLS.CA.File != "" {
		server.MustCreateCA()
	}

	// Cert
	if config.Globals.TLS.Cert.File != "" {
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
