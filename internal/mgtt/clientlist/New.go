package clientlist

import "sync"

var listMutex sync.Mutex
var list map[string]Client = make(map[string]Client)

// Init will init the list of clients
func init() {
	if list == nil {
		list = make(map[string]Client)
	}
}
