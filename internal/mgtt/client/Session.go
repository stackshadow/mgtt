package client

func (client *MgttClient) CleanSessionSet(clean bool) {
	client.cleanSession = clean
}

func (client *MgttClient) CleanSessionGet() (clean bool) {
	return client.cleanSession
}
