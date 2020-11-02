package cli

import (
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/broker"
	"gitlab.com/mgtt/plugins/acl"
	"gitlab.com/mgtt/plugins/auth"
)

// CmdServe represents the flag which create certs
type CmdServe struct {
	URL        string `help:"Set the url where the broker is listening"  env:"URL" default:"tcp://0.0.0.0:8883"`
	TLS        bool   `help:"Enable TLS"  env:"TLS" default:"true"`
	CAFile     string `help:"The ca to use for TLS"  env:"CA" default:"tls/ca.crt"`
	CertFile   string `help:"The certificate to use for TLS"  env:"CERT" default:"tls/server.crt"`
	KeyFile    string `help:"The private key to use for TLS"  env:"KEY" default:"tls/server.key"`
	DBFilename string `help:"Filename for retained message-db"  env:"DBFILENAME" default:"messages.db"`
}

// Run will run the command
func (c *CmdServe) Run() (err error) {

	if CLI.Serve.TLS == true {
		// create ca and server-cert
		CLI.CreateCA.Run()
		CLI.CreateCert.Run()
	}

	// set the broker-connection-timeout
	broker.ConnectTimeout = CLI.ConnectTimeout

	// create the broker
	newbroker, err := broker.New()
	if err != nil {
		log.Error().Err(err).Send()
	}

	// register plugins
	if strings.Contains(CLI.Plugins, "auth") == true {
		auth.LocalInit(CLI.ConfigPath)
	}
	if strings.Contains(CLI.Plugins, "acl") == true {
		acl.LocalInit(CLI.ConfigPath)
	}

	err = newbroker.Serve(
		broker.Config{
			URL:        c.URL,
			TLS:        c.TLS,
			CAFile:     c.CAFile,
			CertFile:   c.CertFile,
			KeyFile:    c.KeyFile,
			DBFilename: c.DBFilename,
		},
	)
	if err != nil {
		log.Error().Err(err).Send()
	}
	return
}
