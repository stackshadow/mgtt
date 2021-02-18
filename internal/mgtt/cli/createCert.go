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
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CmdCreateCert represents the flag which create certs
type CmdCreateCert struct {
	CAFile   string `help:"The ca to use for TLS, set this to '' or use self-signed to create an self-signed-certificate"  env:"CAFILE" default:"tls/ca.crt"`
	CertFile string `help:"The certificate-file that will be created ( if not exist )"  env:"CERTFILE" default:"tls/server.crt"`
	KeyFile  string `help:"The private-key-file that will be created ( if not exist )"  env:"KEYFILE" default:"tls/server.key"`

	Organization  string `help:"Organisation of the ca" env:"ORGANIZATION" default:"FeelGood Inc."`
	Country       string `help:"Country-Code" env:"COUNTRY" default:"DE"`
	Province      string `help:"Province" env:"PROVINCE" default:"Local"`
	Locality      string `help:"Locality (City)" env:"LOCALITY" default:"Berlin"`
	StreetAddress string `help:"Adress" env:"STREETADDRESS" default:"Corner 42"`
	PostalCode    string `help:"PostalCode" env:"POSTALCODE" default:"030423"`

	SelfSigned bool `help:"Create self signed certificate" env:"SELFSIGNED" default:"false"`
}

// Run will run the command
func (c *CmdCreateCert) Run() (err error) {

	baseDirName := filepath.Dir(c.CAFile)
	certificateFileName := c.CertFile
	certificatePrivKeyFileName := c.KeyFile

	if c.SelfSigned == true {
		c.CAFile = ""
	}

	// check if the files already exist
	if _, statErr := os.Stat(certificateFileName); !os.IsNotExist(statErr) {
		log.Info().Str("Certificate", certificateFileName).Msg("Already exist, no need to create it")
		return
	}
	if _, statErr := os.Stat(certificatePrivKeyFileName); !os.IsNotExist(statErr) {
		log.Info().Str("Private-Key", certificatePrivKeyFileName).Msg("Already exist, no need to create it")
		return
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
	}

	if c.CAFile != "" {
		cert.KeyUsage = x509.KeyUsageDigitalSignature
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	} else {
		cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	var certPrivKey *rsa.PrivateKey
	certPrivKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var certBytes []byte

	if c.CAFile != "" {
		baseFileName := filepath.Base(strings.TrimSuffix(c.CAFile, path.Ext(c.CAFile)))
		caCertFileName := baseDirName + "/" + baseFileName + ".crt"
		caPrivKeyFileName := baseDirName + "/" + baseFileName + ".key"

		// Load CA
		var catls tls.Certificate
		catls, err = tls.LoadX509KeyPair(caCertFileName, caPrivKeyFileName)
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
		if err != nil {
			return err
		}
	} else {

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
		if err != nil {
			return err
		}

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
	os.Mkdir(filepath.Dir(baseDirName), 0700)

	// write file
	err = ioutil.WriteFile(certificateFileName, certPEM.Bytes(), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write certificate")
	}

	// write file
	err = ioutil.WriteFile(certificatePrivKeyFileName, certPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write key")
	}

	return
}
