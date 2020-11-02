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
	CAFile     string `help:"The ca to use for TLS. If not set, TLS is disabled"  env:"CA" default:"tls/ca.crt.pem"`
	CertFile   string `help:"The certificate to use for TLS. If not set, TLS is disabled"  env:"CERT" default:"tls/server.crt.pem"`
	KeyFile    string `help:"The private key to use for TLS"  env:"KEY" default:"tls/server.key.pem"`
	DBFilename string `help:"Filename for retained message-db"  env:"DBFILENAME" default:"messages.db"`
}

// Run will run the command
func (c *CmdServe) Run() (err error) {

	broker.ConnectTimeout = CLI.ConnectTimeout

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
