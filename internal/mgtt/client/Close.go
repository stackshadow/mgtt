package client

// Close will close an network connection
func (client *MgttClient) Close() (err error) {

	// resend-clients we don't close
	if client.id == "resend" {
		return nil
	}

	err = client.connection.Close()

	// close the loop
	client.packetSendLoopExit <- true

	client.Connected = false
	return
}
