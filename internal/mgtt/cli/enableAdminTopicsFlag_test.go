package cli

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
)

func TestAdminTopicFlag(t *testing.T) {

	parser, err := kong.New(&cliData)
	if err != nil {
		panic(err)
	}
	_, err = parser.Parse([]string{
		"--enable-admin-topics",
		"serve",
	})
	parser.FatalIfErrorf(err)

	// check if the environment variable is set to true
	if os.Getenv("ENABLE_ADMIN_TOPICS") != "true" {
		t.FailNow()
	}

}
