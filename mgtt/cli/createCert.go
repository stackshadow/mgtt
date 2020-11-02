package cli

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
)

// CmdCreateCert represents the flag which create certs
type CmdCreateCert struct {
	OutputDirectory string `help:"Ouput directory" default:"tls"`
	Organization    string `help:"Organisation of the ca" default:"FeelGood Inc."`
	Country         string `help:"Country-Code" default:"DE"`
	Province        string `help:"Province" default:"Local"`
	Locality        string `help:"Locality (City)" default:"Berlin"`
	StreetAddress   string `help:"Adress" default:"Corner 42"`
	PostalCode      string `help:"PostalCode" default:"030423"`
	Name            string `help:"Name of new cert" default:"client"`
}

// Run will run the command
func (c *CmdCreateCert) Run() (err error) {

	// cert-data
	var caCert *x509.Certificate

	// Load CA
	var catls tls.Certificate
	catls, err = tls.LoadX509KeyPair(c.OutputDirectory+"/ca.crt.pem", c.OutputDirectory+"/ca.key.pem")
	if err != nil {
		panic(err)
	}
	caCert, err = x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	// caCert, err = tls.LoadX509KeyPair(c.OutputDirectory+"/ca.crt.pem", c.OutputDirectory+"/ca.key.pem")

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    "MGTT - Client",
			Organization:  []string{c.Organization},
			Country:       []string{c.Country},
			Province:      []string{c.Province},
			Locality:      []string{c.Locality},
			StreetAddress: []string{c.StreetAddress},
			PostalCode:    []string{c.PostalCode},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	var certPrivKey *rsa.PrivateKey
	certPrivKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var certBytes []byte
	certBytes, err = x509.CreateCertificate(rand.Reader, cert, caCert, &certPrivKey.PublicKey, catls.PrivateKey)
	if err != nil {
		return err
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

	// create subdirectory
	os.Mkdir(filepath.Dir(c.OutputDirectory), 0700)

	// write file
	err = ioutil.WriteFile(c.OutputDirectory+"/"+c.Name+".crt.pem", certPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write certificate")
	}

	// write file
	err = ioutil.WriteFile(c.OutputDirectory+"/"+c.Name+".key.pem", certPrivKeyPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write key")
	}

	return
}
