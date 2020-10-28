package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) loopHandleClientPackets(connectedClient *client.MgttClient, packet packets.ControlPacket) {
	for {
		event := <-broker.clientEvents

		log.Debug().
			Uint16("mid", event.packet.Details().MessageID).
			Uint8("Qos", event.packet.Details().Qos).
			Str("packet", event.packet.String()).
			Msg("Handle packet")

		var err error = nil

		// CONNACK-Packet
		switch recvPacket := event.packet.(type) {
		case *packets.ConnackPacket:
			err = broker.handleConackPacket(connectedClient, recvPacket)
		}

		// check if client connects correctly
		if event.client != nil {
			if event.client.Connected == false {
				log.Error().Msg("Client not send an CONECT-Packet")
				continue
			}
		}

		switch packet := event.packet.(type) {

		case *packets.PublishPacket:

			if packet.Qos == client.SubackQoS1 {
				event.client.SendPuback(packet.MessageID)
			}

		}

		if err != nil {
			log.Error().Err(err).Send()
		}

	}
}
