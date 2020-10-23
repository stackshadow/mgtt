package client

import "github.com/eclipse/paho.mqtt.golang/packets"

// SendConnect will send an CON-Package to the remote
func (client *MgttClient) SendConnect(username, password, clientid string) (err error) {

	// construct the package
	connect := packets.NewControlPacket(packets.Connect).(*packets.ConnectPacket)
	connect.Username = username
	connect.Password = []byte(password)
	connect.ClientIdentifier = clientid

	// send it
	err = connect.Write(client.connection)

	return
}
