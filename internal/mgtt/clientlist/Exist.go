package clientlist

// Exist return true if clientID exist in the list
func Exist(clientID string) (exist bool) {
	// mutex
	listMutex.Lock()
	defer listMutex.Unlock()

	_, exist = list[clientID]
	return
}
