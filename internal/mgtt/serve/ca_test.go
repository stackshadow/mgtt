package server

import (
	"testing"

	"gitlab.com/mgtt/internal/mgtt/config"
)

func TestCreateOfCA(t *testing.T) {

	config.Load("inttest.yaml")
	config.Values.Level = "debug"
	config.Apply()

	MustCreateCA()
}
