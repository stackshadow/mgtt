package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPingreq will send an PINGREQ-Package
func (client *MgttClient) SendPingreq() {

	// construct the package
	packet := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)

	log.Debug().
		Str("cid", client.ID()).
		Msg("Send PINGREQ")

	// queue packet
	client.sendPackets <- packet

	return
}
