package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/broker"
	"gitlab.com/mgtt/cli"
)

func init() {
	// ########################## Command line parse ##########################
	kong.Parse(&cli.CLI,
		kong.Name("mgtt"),
		kong.Description("Message Go Telemetry Transport"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	cli.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped

	// we check if key is set and exist
	// if not exist, we create a key for you :)
	if cli.CLI.KeyFile != "" {
		if _, err := os.Stat(cli.CLI.KeyFile); err != nil {

			// open keyfile
			keyOut, err := os.OpenFile(cli.CLI.KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to open key for writing")
			}

			// create an private key
			_, priv, err := ed25519.GenerateKey(rand.Reader)

			// marshall it
			privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to marshal private key")
			}

			// write it
			if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
				log.Fatal().Err(err).Msg("Failed to write data to key")
			}

			// close it
			if err := keyOut.Close(); err != nil {
				log.Fatal().Err(err).Msg("Failed to write data to key")
			}

		}
	}

	if cli.CLI.CertFile != "" {
		if _, err := os.Stat(cli.CLI.CertFile); err != nil {

			// open keyfile
			privKeyData, err := ioutil.ReadFile(cli.CLI.KeyFile)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to open key for writing")
			}

			block, _ := pem.Decode(privKeyData)

			privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to open key for writing")
			}

			privEDKey := privKey.(ed25519.PrivateKey)

			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

			template := x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Organization: []string{"mgtt local"},
				},

				NotBefore: time.Now(),
				NotAfter:  time.Now().Add(time.Hour * 24 * 356),
				KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,

				ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

				BasicConstraintsValid: true,
				IsCA:                  true,
			}

			certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privEDKey.Public(), privEDKey)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create certificate")
			}

			// write it
			certOut := bytes.NewBufferString("")
			if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
				log.Fatal().Err(err).Msg("Failed to write data to key")
			}

			err = ioutil.WriteFile(cli.CLI.CertFile, certOut.Bytes(), 0777)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create certificate")
			}
		}
	}

}

func main() {

	newbroker, err := broker.New()
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = newbroker.Serve(
		broker.Config{
			URL:      cli.CLI.URL,
			CertFile: cli.CLI.CertFile,
			KeyFile:  cli.CLI.KeyFile,
		},
	)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
