package client

import (
	"github.com/rs/zerolog/log"
)

// PacketSendLoop start an loop ( blocking ) that wait for client.sendPackets and send it over the network connection
//
// If the connection is closed or an error occured, the loop will end
func (c *MgttClient) packetSendLoop() {

	c.packetSendLoopRunning = true

loop:
	for {

		select {
		case <-c.packetSendLoopExit:
			break loop
		case packetToSend := <-c.sendPackets:

			// connected ?
			if c.connection == nil {
				log.Warn().
					Str("client", c.ID()).
					Msg("can not send packet, connection is closed")
				break
			}

			// send it
			err := packetToSend.Write(c.connection)
			if err != nil {
				log.Error().
					Str("client", c.ID()).
					Err(err).Send()
				break
			}

			log.Debug().
				Str("client", c.ID()).
				Str("packet", packetToSend.String()).
				Msg("packet send")

		}

	}

	c.packetSendLoopRunning = false
	return
}
