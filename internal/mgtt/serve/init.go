package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

// MustInit init the listener or panics
func (l *Listener) MustInit(CA, Cert, Key string) {
	var err error

	if Cert == "" || Key == "" {
		l.listener, err = net.Listen("tcp", l.address)
		utils.PanicOnErr(err)

		log.Info().Str("create", l.address).
			Str("ca", CA).
			Str("cert", Cert).
			Str("key", Key).
			Msg("Listening non-tls")
	}
	if Cert != "" && Key != "" {
		TLSConfig := mustTLSConfig(CA, Cert, Key)
		l.listener, err = tls.Listen("tcp", l.address, TLSConfig)
		utils.PanicOnErr(err)

		log.Info().Str("create", l.address).
			Str("ca", CA).
			Str("cert", Cert).
			Str("key", Key).
			Msg("Listening")
	}

}

// getTLSConfig return the current tls-config or err if an error occured
func mustTLSConfig(CA, Cert, Key string) (cfg *tls.Config) {

	var err error
	var cert tls.Certificate
	var clientCAs *x509.CertPool
	var clientCABytes []byte

	cert, err = tls.LoadX509KeyPair(Cert, Key)
	utils.PanicOnErr(err)

	// Get the SystemCertPool, continue with an empty pool on error
	clientCAs, err = x509.SystemCertPool()
	if clientCAs == nil {
		clientCAs = x509.NewCertPool()
	}

	// CA was set = mTLS
	if CA != "" {
		// Read in the cert file

		clientCABytes, err = ioutil.ReadFile(CA)
		utils.PanicOnErr(err)

		// Append our cert to the system pool
		if ok := clientCAs.AppendCertsFromPEM(clientCABytes); !ok {
			log.Warn().
				Msg("No certs appended, using system certs only")
		}

	}

	cfg = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		ClientCAs:          clientCAs,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		},
	}

	if CA != "" {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		log.Info().Msg("using mTLS")
	} else {
		cfg.ClientAuth = tls.NoClientCert
		log.Info().Msg("using self-signed-certificate")
	}

	return
}
