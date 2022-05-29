package server

import (
	"testing"

	"gitlab.com/mgtt/internal/mgtt/config"
)

func TestCreateOfCA(t *testing.T) {

	config.MustLoadFromFile("inttest.yaml")
	config.Globals.Level = "debug"
	config.Globals.TLS.CA.File = "tls_ca.crt"
	config.ApplyLog()

	MustCreateCA()
}
