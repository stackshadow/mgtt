package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendConnect will send an CON-Package to the remote
func (client *MgttClient) SendConnect(username, password, clientid string) (err error) {

	// construct the package
	connect := packets.NewControlPacket(packets.Connect).(*packets.ConnectPacket)
	connect.Username = username
	connect.Password = []byte(password)
	connect.ClientIdentifier = clientid

	log.Debug().
		Str("cid", clientid).
		Str("username", username).
		Msg("Send CONNECKT")

	// send it
	err = connect.Write(client.connection)

	return
}
