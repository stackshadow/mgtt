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
	CAFile     string `help:"The ca to use for TLS. Set it to an empty string if you would like to use an self-signed-certificate" env:"CAFILE" default:"tls/ca.crt"`
	CertFile   string `help:"The certificate to use for TLS" env:"CERTFILE" default:"tls/server.crt"`
	KeyFile    string `help:"The private key to use for TLS" env:"KEYFILE" default:"tls/server.key"`
	DBFilename string `help:"Filename for retained message-db" env:"DBFILENAME" default:"messages.db"`

	SelfSigned bool `help:"Use self-signed-certificate and ignore CAFile" env:"SELFSIGNED" default:"false"`
}

// Run will run the command
func (c *CmdServe) Run(params Parameter) (err error) {

	// did we use TLS ?
	if params.Serve.TLS == true {

		if c.SelfSigned == true {
			c.CAFile = ""
		}

		// create ca
		if c.CAFile != "" {
			params.CreateCA.CAFile = c.CAFile
			params.CreateCA.Run()
		}

		// create certificate if not exist
		params.CreateCert.CAFile = c.CAFile
		params.CreateCert.CertFile = c.CertFile
		params.CreateCert.KeyFile = c.KeyFile
		params.CreateCert.SelfSigned = c.SelfSigned
		params.CreateCert.Run()
	}

	// set the broker-connection-timeout
	broker.ConnectTimeout = params.ConnectTimeout

	// create the broker
	newbroker, err := broker.New()
	if err != nil {
		log.Error().Err(err).Send()
	}

	// register plugins
	if strings.Contains(params.Plugins, "auth") == true {
		auth.LocalInit(params.ConfigPath)
	}
	if strings.Contains(params.Plugins, "acl") == true {
		acl.LocalInit(params.ConfigPath)
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

	var done chan bool
	done, err = newbroker.Serve(newBrokerConfig)
	if err != nil {
		log.Error().Err(err).Send()
	}
	<-done
	return
}
