package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// respond
// 0x00 	Connection Accepted
// 0x01 	Connection Refused: unacceptable protocol version
// 0x02 	Connection Refused: identifier rejected
// 0x03 	Connection Refused: server unavailable
// 0x04 	Connection Refused: bad user name or password
// 0x05 	Connection Refused: not authorized
const (
	ConnackAccepted            = 0x00
	ConnackUnacceptable        = 0x01
	ConnackIDRejected          = 0x02
	ConnackServerUnavailable   = 0x03
	ConnackBadUsernamePassword = 0x04
	ConnackUnauthorized        = 0x05
)

// SendConnack will send an CONACK-Package to the client
func (client *MgttClient) SendConnack(ReturnCode byte, SessionPresent bool) (err error) {

	// construct the package
	packet := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)
	packet.ReturnCode = ReturnCode
	packet.SessionPresent = SessionPresent

	log.Debug().
		Str("cid", client.ID()).
		Uint8("return code", ReturnCode).
		Msg("Send CONACK")

	// queue packet
	client.sendPackets <- packet

	return
}
