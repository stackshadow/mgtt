package config

import (
	"testing"

	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

func TestSaveConfig(t *testing.T) {

	testConfig := `
level: debug
`

	// load the config
	LoadFromByte([]byte(testConfig))

	// change the config
	Values.URL = "ChangedURL"

	// save it to a file
	fileName = "./integrationtest_saveconfig.yml"
	MustSave()

	// clear it
	Values.URL = ""

	// load it again
	MustLoad(fileName)

	// check if
	if Values.URL != "ChangedURL" {
		utils.PanicWithString("URL can not be changed")
	}
}
