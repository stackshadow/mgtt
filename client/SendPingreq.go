package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPingreq will send an PINGREQ-Package
func (client *MgttClient) SendPingreq() (err error) {

	// construct the package
	pingReq := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)

	log.Debug().
		Str("cid", client.ID()).
		Msg("Send PINGREQ")

	// send it
	err = pingReq.Write(client.connection)

	return
}
