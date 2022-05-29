package acl

import (
	"sync"
)

type pluginConfigStruct struct {
	Enable bool                  `yaml:"enable"`
	Rules  map[string][]aclEntry `yaml:"rules"`
}

var mutex sync.Mutex
var pluginConfig *pluginConfigStruct = &pluginConfigStruct{}
