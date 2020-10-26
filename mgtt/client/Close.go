package client

// Close will close an network connection
func (client *MgttClient) Close() (err error) {

	err = client.connection.Close()

	client.Connected = false
	return
}
