package client

import "github.com/eclipse/paho.mqtt.golang/packets"

// SendPingresp will send an PINGRESP-Package
func (client *MgttClient) SendPingresp() {

	pingResp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
	pingResp.Write(client.connection)

	return
}
