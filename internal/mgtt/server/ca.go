package server

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
	"gitlab.com/mgtt/internal/mgtt/config"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

func MustCreateCA() {

	// ca should not created
	if config.Globals.TLS.CA.File == "" {
		log.Debug().Msg("not create ca")
		return
	}

	var err error

	// ca already created
	var CAFileAbsolute string
	CAFileAbsolute, err = filepath.Abs(config.Globals.TLS.CA.File)
	utils.PanicOnErr(err)

	// already exist
	if _, statErr := os.Stat(CAFileAbsolute); !os.IsNotExist(statErr) {
		log.Info().Str("CA-Certificate", CAFileAbsolute).Msg("Already exist, no need to create it")
		return
	}

	// filepath.Ext() filepath.Base(CLI.Serve.CAFile)
	caPrivKeyFileName := CAFileAbsolute + ".key"

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    "MGTT CA",
			Organization:  []string{config.Globals.TLS.CA.Organization},
			Country:       []string{config.Globals.TLS.CA.Country},
			Province:      []string{config.Globals.TLS.CA.Province},
			Locality:      []string{config.Globals.TLS.CA.Locality},
			StreetAddress: []string{config.Globals.TLS.CA.StreetAddress},
			PostalCode:    []string{config.Globals.TLS.CA.PostalCode},
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
	utils.PanicOnErr(err)

	var caBytes []byte
	caBytes, err = x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	utils.PanicOnErr(err)

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

	// write file
	err = ioutil.WriteFile(CAFileAbsolute, caPEM.Bytes(), 0600)
	utils.PanicOnErr(err)

	// write file
	err = ioutil.WriteFile(caPrivKeyFileName, caPrivKeyPEM.Bytes(), 0600)
	utils.PanicOnErr(err)

}
