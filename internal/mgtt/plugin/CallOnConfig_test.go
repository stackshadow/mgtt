package plugin

import (
	"testing"

	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

type testPluginConfig struct {
	Enable   bool   `yaml:"enable"`
	TestData string `yaml:"data"`
}

// create the plugin
var testPlugin = V1{
	OnPluginConfig: func(pluginData *interface{}) (configChanged bool) {

		// get data as struct
		var testPluginConfigData testPluginConfig
		tempData, err := yaml.Marshal(pluginData)
		utils.PanicOnErr(err)
		yaml.Unmarshal(tempData, &testPluginConfigData)

		// we check it
		if testPluginConfigData.Enable != true {
			utils.PanicWithString("Can not parse config")
		}

		// we alter the config
		testPluginConfigData.TestData = "altered"

		// push it back
		// this is needed, so that the plugin-system can store your config !
		tempData, err = yaml.Marshal(testPluginConfigData)
		utils.PanicOnErr(err)
		yaml.Unmarshal(tempData, pluginData)

		// we inform the plugin-system that it should store the config
		configChanged = true

		return
	},
}

func TestPluginConfig(t *testing.T) {

	// we prepare the plugins
	var PluginConfigs map[string]interface{} = make(map[string]interface{})
	PluginConfigs["test"] = testPluginConfig{Enable: true}

	// register the plugin
	Register("test", &testPlugin)

	// call all plugins
	CallOnPluginConfig(PluginConfigs)

	// now we start to check if the plugin altered the config
	var testPluginConfigData testPluginConfig
	tempData, err := yaml.Marshal(PluginConfigs["test"])
	utils.PanicOnErr(err)
	yaml.Unmarshal(tempData, &testPluginConfigData)

	if testPluginConfigData.TestData != "altered" {
		utils.PanicWithString("Testdata was not corret changed by plugin")
	}
}
