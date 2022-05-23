package config

import "gopkg.in/yaml.v2"

func AlterPluginConfig(pluginName string, pluginData interface{}) {

	plugins := valuesRawMap["plugins"].(map[interface{}]interface{})
	plugins[pluginName] = pluginData
	valuesRawMap["plugins"] = plugins

}

func GetPluginData(pluginName string, pluginData interface{}) (err error) {

	plugins := valuesRawMap["plugins"].(map[interface{}]interface{})
	if plugins != nil {
		pluginInterface := plugins[pluginName]
		if pluginInterface != nil {
			var tempBytes []byte
			tempBytes, err = yaml.Marshal(pluginInterface)
			if err == nil {
				err = yaml.Unmarshal(tempBytes, pluginData)
			}
		}
	}

	return
}
