package cli

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CmdCreateCA reflect the create-ca command
type CmdCreateCA struct {
	CAFile string `help:"The ca to use for TLS"  env:"CA" default:"tls/ca.crt"`

	Organization  string `help:"Organisation of the ca" default:"FeelGood Inc."`
	Country       string `help:"Country-Code" default:"DE"`
	Province      string `help:"Province" default:"Local"`
	Locality      string `help:"Locality (City)" default:"Berlin"`
	StreetAddress string `help:"Adress" default:"Corner 42"`
	PostalCode    string `help:"PostalCode" default:"030423"`
}

// Run will create a new CA
func (c *CmdCreateCA) Run() (err error) {

	// filepath.Ext() filepath.Base(CLI.Serve.CAFile)
	baseDirName := filepath.Dir(c.CAFile)
	baseFileName := filepath.Base(strings.TrimSuffix(c.CAFile, path.Ext(c.CAFile)))
	caCertFileName := baseDirName + "/" + baseFileName + ".crt"
	caPrivKeyFileName := baseDirName + "/" + baseFileName + ".key"

	// create directory if needed
	os.MkdirAll(baseDirName, 0777)

	// check if the files already exist
	if _, statErr := os.Stat(caCertFileName); !os.IsNotExist(statErr) {
		log.Info().Str("CA-Certificate", caCertFileName).Msg("Already exist, no need to create it")
		return
	}
	if _, statErr := os.Stat(caPrivKeyFileName); !os.IsNotExist(statErr) {
		log.Info().Str("CA-Private-Key", caPrivKeyFileName).Msg("Already exist, no need to create it")
		return
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    "MGTT CA",
			Organization:  []string{c.Organization},
			Country:       []string{c.Country},
			Province:      []string{c.Province},
			Locality:      []string{c.Locality},
			StreetAddress: []string{c.StreetAddress},
			PostalCode:    []string{c.PostalCode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	var caPrivKey *rsa.PrivateKey
	caPrivKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var caBytes []byte
	caBytes, err = x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	var privBytes []byte
	privBytes, err = x509.MarshalPKCS8PrivateKey(caPrivKey)
	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	// create subdirectory
	os.Mkdir(filepath.Dir(baseDirName), 0700)

	// write file
	err = ioutil.WriteFile(caCertFileName, caPEM.Bytes(), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write certificate")
	}

	// write file
	err = ioutil.WriteFile(caPrivKeyFileName, caPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write key")
	}

	return nil
}
