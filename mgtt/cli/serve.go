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
	CAFile     string `help:"The ca to use for TLS. Set it to an empty string if you would like to use an self-signed-certificate"  env:"CA" default:"tls/ca.crt"`
	CertFile   string `help:"The certificate to use for TLS"  env:"CERT" default:"tls/server.crt"`
	KeyFile    string `help:"The private key to use for TLS"  env:"KEY" default:"tls/server.key"`
	DBFilename string `help:"Filename for retained message-db"  env:"DBFILENAME" default:"messages.db"`

	SelfSigned bool `help:"Use self-signed-certificate and ignore CAFile" default:"false"`
}

// Run will run the command
func (c *CmdServe) Run() (err error) {

	if CLI.Serve.TLS == true {
		// create ca and server-cert
		CLI.CreateCA.CAFile = c.CAFile
		CLI.CreateCA.Run()

		CLI.CreateCert.CAFile = c.CAFile
		CLI.CreateCert.CertFile = c.CertFile
		CLI.CreateCert.KeyFile = c.KeyFile
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

	newBrokerConfig := broker.Config{
		URL:        c.URL,
		TLS:        c.TLS,
		CAFile:     c.CAFile,
		CertFile:   c.CertFile,
		KeyFile:    c.KeyFile,
		DBFilename: c.DBFilename,
	}

	if c.SelfSigned == true {
		c.CAFile = ""
	}

	err = newbroker.Serve(newBrokerConfig)
	if err != nil {
		log.Error().Err(err).Send()
	}
	return
}
