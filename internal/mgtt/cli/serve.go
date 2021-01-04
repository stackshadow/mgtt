package cli

import (
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt-plugins/acl"
	"gitlab.com/mgtt/internal/mgtt-plugins/auth"
	"gitlab.com/mgtt/internal/mgtt/broker"
)

// CmdServe represents the flag which create certs
type CmdServe struct {
	URL        string `help:"Set the url where the broker is listening" env:"URL" default:"tcp://0.0.0.0:8883"`
	TLS        bool   `help:"Enable TLS" env:"TLS" default:"true"`
	CAFile     string `help:"The ca to use for TLS. Set it to an empty string if you would like to use an self-signed-certificate" env:"CA" default:"tls/ca.crt"`
	CertFile   string `help:"The certificate to use for TLS" env:"CERT" default:"tls/server.crt"`
	KeyFile    string `help:"The private key to use for TLS" env:"KEY" default:"tls/server.key"`
	DBFilename string `help:"Filename for retained message-db" env:"DBFILENAME" default:"messages.db"`

	SelfSigned bool `help:"Use self-signed-certificate and ignore CAFile" env:"SELFSIGNED" default:"false"`
}

// Run will run the command
func (c *CmdServe) Run(cliCurrent CLIType) (err error) {

	// did we use TLS ?
	if cliCurrent.Serve.TLS == true {

		if c.SelfSigned == true {
			c.CAFile = ""
		}

		// create ca
		if c.CAFile != "" {
			cliCurrent.CreateCA.CAFile = c.CAFile
			cliCurrent.CreateCA.Run()
		}

		// create certificate if not exist
		cliCurrent.CreateCert.CAFile = c.CAFile
		cliCurrent.CreateCert.CertFile = c.CertFile
		cliCurrent.CreateCert.KeyFile = c.KeyFile
		cliCurrent.CreateCert.SelfSigned = c.SelfSigned
		cliCurrent.CreateCert.Run()
	}

	// set the broker-connection-timeout
	broker.ConnectTimeout = cliCurrent.ConnectTimeout

	// create the broker
	newbroker, err := broker.New()
	if err != nil {
		log.Error().Err(err).Send()
	}

	// register plugins
	if strings.Contains(cliCurrent.Plugins, "auth") == true {
		auth.LocalInit(cliCurrent.ConfigPath)
	}
	if strings.Contains(cliCurrent.Plugins, "acl") == true {
		acl.LocalInit(cliCurrent.ConfigPath)
	}

	newBrokerConfig := broker.Config{
		Version:    Version,
		URL:        c.URL,
		TLS:        c.TLS,
		CAFile:     c.CAFile,
		CertFile:   c.CertFile,
		KeyFile:    c.KeyFile,
		DBFilename: c.DBFilename,
	}

	err = newbroker.Serve(newBrokerConfig)
	if err != nil {
		log.Error().Err(err).Send()
	}
	return
}
