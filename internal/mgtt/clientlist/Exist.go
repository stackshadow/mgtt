package clientlist

// Exist return true if clientID exist in the list
func Exist(clientID string) (exist bool) {
	_, exist = list[clientID]
	return
}
