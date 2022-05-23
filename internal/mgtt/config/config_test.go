package config

import (
	"testing"

	"gitlab.com/mgtt/internal/mgtt/plugin"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v3"
)

type testPluginConfig struct {
	Enable   bool   `yaml:"enable"`
	TestData string `yaml:"data"`
}

var testPluginConfigData testPluginConfig

// create the plugin
var testPlugin = plugin.V1{
	OnPluginConfig: func(pluginData []byte) (configChanged bool) {

		err := yaml.Unmarshal(pluginData, &testPluginConfigData)
		utils.PanicOnErr(err)

		if testPluginConfigData.Enable != true {
			utils.PanicWithString("Can not parse config")
		}

		// we alter the config
		testPluginConfigData.TestData = "altered"

		AlterPluginConfig("test", testPluginConfigData)

		return
	},
}

func TestPluginConfig(t *testing.T) {
	var err error

	testConfig := `
level: debug
plugins:
  test:
    enable: true
    data: integrationtest
`

	// register it
	plugin.Register("test", &testPlugin)

	// load the config
	LoadFromByte([]byte(testConfig))

	// parse it again and see it config was altered
	var alteredConfig testPluginConfig

	GetPluginData("test", &alteredConfig)
	utils.PanicOnErr(err)

	if alteredConfig.TestData != "altered" {
		utils.PanicWithString("Data was not altered")
	}
}
