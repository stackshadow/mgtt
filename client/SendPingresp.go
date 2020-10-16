package client

import "github.com/eclipse/paho.mqtt.golang/packets"

// SendPingresp will send an PINGRESP-Package
func (client *MgttClient) SendPingresp() (err error) {

	// construct the package
	pingResp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)

	// send it
	err = pingResp.Write(client.connection)

	return
}
