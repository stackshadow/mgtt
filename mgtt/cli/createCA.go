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
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// DebugFlag represents the flag which enable debugging
type CmdCreateCA struct {
	OutputDirectory string `arg help:"Ouput directory" default:"tls"`
	Organization    string `arg help:"Organisation of the ca" default:"FeelGood Inc."`
	Country         string `arg help:"Country-Code" default:"DE"`
	Province        string `arg help:"Province" default:"Local"`
	Locality        string `arg help:"Locality (City)" default:"Berlin"`
	StreetAddress   string `arg help:"Adress" default:"Corner 42"`
	PostalCode      string `arg help:"PostalCode" default:"030423"`
}

func (c *CmdCreateCA) Run() (err error) {

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
	os.Mkdir(filepath.Dir(c.OutputDirectory), 0700)

	// write file
	err = ioutil.WriteFile(c.OutputDirectory+"/ca.crt.pem", caPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write certificate")
	}

	// write file
	err = ioutil.WriteFile(c.OutputDirectory+"/ca.key.pem", caPrivKeyPEM.Bytes(), 0777)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write key")
	}

	return nil
}
