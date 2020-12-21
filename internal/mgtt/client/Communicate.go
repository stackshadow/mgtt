package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// Communicate will read all packets to the internal channel
// this is an NON-BLOCKING function that start a go-routine
func (client *MgttClient) Communicate() {

	go func() {

		var err error

		for {

			var recvdPacket packets.ControlPacket
			recvdPacket, err = client.readPacket()
			if err != nil {
				close(client.recvPackets)
				break
			}

			client.recvPackets <- recvdPacket
		}

		log.Error().Err(err).Send()
		//client.recvPackets <- nil
	}()

	go func() {
		packetToSend := <-client.sendPackets
		if client.connection == nil {
			log.Warn().
				Str("client", client.ID()).
				Msg("can not send packet, connection is closed")
			return
		}
		err := packetToSend.Write(client.connection)
		if err != nil {
			log.Error().
				Str("client", client.ID()).
				Err(err).Send()
			return
		}

		log.Debug().
			Str("client", client.ID()).
			Str("packet", packetToSend.String()).
			Msg("packet send")
	}()

	return
}
