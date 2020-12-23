package clientlist

import "gitlab.com/mgtt/internal/mgtt/client"

var list map[string]*client.MgttClient

// Init will init the list of clients
func Init() {
	if list == nil {
		list = make(map[string]*client.MgttClient)
	}
}
