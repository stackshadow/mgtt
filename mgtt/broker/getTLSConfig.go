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

	cfg = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		ClientCAs:          clientCAs,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	}

	return
}
