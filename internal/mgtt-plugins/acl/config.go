package acl

import (
	"sync"
)

type pluginConfig struct {
	Plugins struct {
		ACL struct {
			Rules map[string][]aclEntry `yaml:"rules"`
		} `yaml:"acl"`
	} `yaml:"plugins"`
}

var mutex sync.Mutex
var config *pluginConfig = &pluginConfig{}
