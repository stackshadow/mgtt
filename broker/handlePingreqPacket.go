package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePingreqPacket(event *client.Event) {
	packet, ok := event.Packet.(*packets.PingreqPacket)
	if ok == false {
		log.Error().Str("clientid", event.Client.ID()).Msg("Expected SubscribePacket")
		return
	}
	log.Debug().
		Str("clientid", event.Client.ID()).
		Str("packet", packet.String()).
		Msg("RCV PingreqPacket")

	event.Client.SendPingresp()
}
