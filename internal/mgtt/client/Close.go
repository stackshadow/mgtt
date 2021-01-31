package client

import "github.com/rs/zerolog/log"

// Close will close an network connection
func (client *MgttClient) Close() (err error) {

	// close network-connection
	err = client.connection.Close()

	// close the loop
	if client.packetSendLoopRunning {
		client.packetSendLoopExit <- true
	} else {
		log.Warn().Str("cid", client.id).Msg("packetSendLoop already closed")
	}

	client.Connected = false
	return
}
