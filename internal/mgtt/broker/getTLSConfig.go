package broker

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

// getTLSConfig return the current tls-config or err if an error occured
func getTLSConfig(config Config) (cfg *tls.Config, err error) {
	var cert tls.Certificate
	var clientCAs *x509.CertPool
	var clientCABytes []byte

	cert, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile)

	// Get the SystemCertPool, continue with an empty pool on error
	if err == nil {
		clientCAs, err = x509.SystemCertPool()
		if clientCAs == nil {
			clientCAs = x509.NewCertPool()
		}
	}

	// no caFile was set
	if config.CAFile != "" {
		// Read in the cert file
		if err == nil {
			clientCABytes, err = ioutil.ReadFile(config.CAFile)
		}

		// Append our cert to the system pool
		if err == nil {
			if ok := clientCAs.AppendCertsFromPEM(clientCABytes); !ok {
				log.Warn().
					Msg("No certs appended, using system certs only")
			}
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

	if config.CAFile != "" {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}

	return
}
