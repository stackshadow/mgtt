package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

func MustCreateCert() {

	var err error

	var CertFileAbsolute string
	CertFileAbsolute, err = filepath.Abs(config.Values.TLS.Cert.File)
	utils.PanicOnErr(err)

	// check if the files already exist
	if _, statErr := os.Stat(CertFileAbsolute); !os.IsNotExist(statErr) {
		log.Info().Str("Certificate", CertFileAbsolute).Msg("Already exist, no need to create it")
		return
	}

	// caCert, err = tls.LoadX509KeyPair(c.OutputDirectory+"/ca.crt.pem", c.OutputDirectory+"/ca.key.pem")

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    "MGTT Cert",
			Organization:  []string{config.Values.TLS.Cert.Organization},
			Country:       []string{config.Values.TLS.Cert.Country},
			Province:      []string{config.Values.TLS.Cert.Province},
			Locality:      []string{config.Values.TLS.Cert.Locality},
			StreetAddress: []string{config.Values.TLS.Cert.StreetAddress},
			PostalCode:    []string{config.Values.TLS.Cert.PostalCode},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
	}

	if config.Values.TLS.CA.File != "" {
		cert.KeyUsage = x509.KeyUsageDigitalSignature
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	} else {
		cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	var certPrivKey *rsa.PrivateKey
	certPrivKey, err = rsa.GenerateKey(rand.Reader, 2048)
	utils.PanicOnErr(err)

	var certBytes []byte

	if config.Values.TLS.CA.File != "" {

		// Load CA
		var catls tls.Certificate
		catls, err = tls.LoadX509KeyPair(config.Values.TLS.CA.File, config.Values.TLS.CA.File+".key")
		if err != nil {
			panic(err)
		}

		// cert-data
		var caCert *x509.Certificate
		caCert, err = x509.ParseCertificate(catls.Certificate[0])
		if err != nil {
			panic(err)
		}

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, caCert, &certPrivKey.PublicKey, catls.PrivateKey)
		utils.PanicOnErr(err)
	} else {

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
		utils.PanicOnErr(err)

	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privBytes []byte
	privBytes, err = x509.MarshalPKCS8PrivateKey(certPrivKey)
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	/*
		serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
		if err != nil {
			return nil, nil, err
		}
	*/

	// write file
	err = ioutil.WriteFile(CertFileAbsolute, certPEM.Bytes(), 0600)
	utils.PanicOnErr(err)

	// write file
	err = ioutil.WriteFile(CertFileAbsolute+".key", certPrivKeyPEM.Bytes(), 0600)
	utils.PanicOnErr(err)

}
