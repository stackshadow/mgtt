package clientlist

func Exist(clientID string) (exist bool) {
	_, exist = list[clientID]
	return
}
