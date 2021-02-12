package clientlist

import "sync"

var listMutex sync.Mutex
var list map[string]Client = make(map[string]Client)

// Init will init the list of clients
func Init() {
	if list == nil {
		list = make(map[string]Client)
	}
}
