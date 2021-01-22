package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) onPacketPinReq(connectedClient *client.MgttClient, packet *packets.PingreqPacket) (err error) {

	connectedClient.SendPingresp()
	return
}
