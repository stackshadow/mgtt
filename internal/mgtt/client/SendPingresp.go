package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPingresp will send an PINGRESP-Package
func (client *MgttClient) SendPingresp() {

	// construct the package
	packet := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)

	log.Debug().
		Str("cid", client.ID()).
		Msg("Send PINGRESP")

	// queue packet
	client.sendPackets <- packet

	return
}
