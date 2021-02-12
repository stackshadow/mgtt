package clientlist

import "errors"

// Move will move an client to an new ID
//
// The client ID will also be set for the new client
func Move(oldID, newID string) (err error) {
	// mutex
	listMutex.Lock()
	defer listMutex.Unlock()

	if currentClient, exist := list[oldID]; exist == true {

		// remove the client
		delete(list, oldID)

		// set an new id
		currentClient.IDSet(newID)

		// and re-add it to the list if possible
		list[currentClient.ID()] = currentClient
	} else {
		err = errors.New("client don't exist")
	}

	return
}
