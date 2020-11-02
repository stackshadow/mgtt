package client

// RemoteAddr return the remoteAddr as string
func (client *MgttClient) RemoteAddr() (remoteAddr string) {
	remoteAddr = client.connection.RemoteAddr().String()
	return
}
