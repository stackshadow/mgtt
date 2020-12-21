package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendConnect will send an CON-Package to the remote
func (client *MgttClient) SendConnect(username, password, clientid string) {

	// construct the package
	packet := packets.NewControlPacket(packets.Connect).(*packets.ConnectPacket)
	packet.Username = username
	packet.Password = []byte(password)
	packet.ClientIdentifier = clientid

	log.Debug().
		Str("cid", clientid).
		Str("username", username).
		Msg("Send CONNECKT")

	// queue packet
	client.sendPackets <- packet

	return
}
