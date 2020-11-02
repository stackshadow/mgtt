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
	CAFile   string `help:"The ca to use for TLS"  env:"CA" default:"tls/ca.crt"`
	CertFile string `help:"The certificate to use for TLS"  env:"CERT" default:"tls/server.crt"`
	KeyFile  string `help:"The private key to use for TLS"  env:"KEY" default:"tls/server.key"`

	Organization  string `help:"Organisation of the ca" default:"FeelGood Inc."`
	Country       string `help:"Country-Code" default:"DE"`
	Province      string `help:"Province" default:"Local"`
	Locality      string `help:"Locality (City)" default:"Berlin"`
	StreetAddress string `help:"Adress" default:"Corner 42"`
	PostalCode    string `help:"PostalCode" default:"030423"`
}

// Run will run the command
func (c *CmdCreateCert) Run() (err error) {

	baseDirName := filepath.Dir(c.CAFile)
	baseFileName := filepath.Base(strings.TrimSuffix(c.CAFile, path.Ext(c.CAFile)))
	caCertFileName := baseDirName + "/" + baseFileName + ".crt"
	caPrivKeyFileName := baseDirName + "/" + baseFileName + ".key"
	certificateFileName := c.CertFile
	certificatePrivKeyFileName := c.KeyFile

	// check if the files already exist
	if _, statErr := os.Stat(certificateFileName); !os.IsNotExist(statErr) {
		log.Info().Str("Certificate", certificateFileName).Msg("Already exist, no need to create it")
		return
	}
	if _, statErr := os.Stat(certificatePrivKeyFileName); !os.IsNotExist(statErr) {
		log.Info().Str("Private-Key", certificatePrivKeyFileName).Msg("Already exist, no need to create it")
		return
	}

	// cert-data
	var caCert *x509.Certificate

	// Load CA
	var catls tls.Certificate
	catls, err = tls.LoadX509KeyPair(caCertFileName, caPrivKeyFileName)
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
	os.Mkdir(filepath.Dir(baseDirName), 0700)

	// write file
	err = ioutil.WriteFile(certificateFileName, certPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write certificate")
	}

	// write file
	err = ioutil.WriteFile(certificatePrivKeyFileName, certPrivKeyPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write key")
	}

	return
}
