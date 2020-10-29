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

	return
}
