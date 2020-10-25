package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// SendPingresp will send an PINGRESP-Package
func (client *MgttClient) SendPingresp() (err error) {

	// construct the package
	pingResp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)

	log.Debug().
		Str("cid", client.ID()).
		Msg("Send PINGRESP")

	// send it
	err = pingResp.Write(client.connection)

	return
}
