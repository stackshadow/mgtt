package server

import (
	"testing"

	"gitlab.com/mgtt/internal/mgtt/config"
)

func TestCreateOfCA(t *testing.T) {

	config.MustLoad("inttest.yaml")
	config.Values.Level = "debug"
	config.Values.TLS.CA.File = "tls_ca.crt"
	config.ApplyLog()

	MustCreateCA()
}
