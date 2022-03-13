package cli

import (
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt-plugins/acl"
	"gitlab.com/mgtt/internal/mgtt-plugins/auth"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
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
	if c.TLS == true {

		if c.SelfSigned == true {
			c.CAFile = ""
		}

		// create ca
		if c.CAFile != "" {
			params.CreateCA.CAFile = c.CAFile
			err = params.CreateCA.Run()
			utils.PanicOnErr(err)
		}

		// create certificate if not exist
		params.CreateCert.CAFile = c.CAFile
		params.CreateCert.CertFile = c.CertFile
		params.CreateCert.KeyFile = c.KeyFile
		params.CreateCert.SelfSigned = c.SelfSigned
		err = params.CreateCert.Run()
		utils.PanicOnErr(err)

	} else {
		log.Warn().Msg("TLS is disabled, this is not a good idea")
	}

	// Broker
	var newbroker *broker.Broker

	// set the broker-connection-timeout
	broker.ConnectTimeout = params.ConnectTimeout

	// create the broker
	newbroker, err = broker.New()
	utils.PanicOnErr(err)

	// register plugins
	if strings.Contains(params.Plugins, "auth") == true {
		auth.LocalInit(params.ConfigPath)
	}
	if strings.Contains(params.Plugins, "acl") == true {
		acl.LocalInit(params.ConfigPath)
	}

	var done chan bool

	newBrokerConfig := broker.Config{
		Version:    Version,
		URL:        c.URL,
		TLS:        c.TLS,
		CAFile:     c.CAFile,
		CertFile:   c.CertFile,
		KeyFile:    c.KeyFile,
		DBFilename: c.DBFilename,
	}

	done, err = newbroker.Serve(newBrokerConfig)
	if err != nil {
		log.Error().Err(err).Send()
	}

	<-done

	return
}
