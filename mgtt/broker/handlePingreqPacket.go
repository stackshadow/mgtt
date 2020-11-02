package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePingreqPacket(connectedClient *client.MgttClient, packet *packets.PingreqPacket) (err error) {

	err = connectedClient.SendPingresp()
	return
}
