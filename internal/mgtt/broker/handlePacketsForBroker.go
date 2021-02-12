package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) handlePacketsForBroker(eventClient *client.MgttClient, eventPacket packets.ControlPacket) (normalClose bool, err error) {

	log.Debug().
		Uint16("mid", eventPacket.Details().MessageID).
		Uint8("Qos", eventPacket.Details().Qos).
		Str("packet", eventPacket.String()).
		Msg("Handle packet")

	switch recvPacket := eventPacket.(type) {
	case *packets.ConnectPacket:
		err = broker.onPacketConnect(eventClient, recvPacket)
		return
	case *packets.DisconnectPacket:
		err = broker.onPacketDisConnect(eventClient)
		normalClose = true
		return
	}

	// check if client connects correctly
	if eventClient.Connected == false {
		log.Error().Str("cid", eventClient.ID()).Msg("Client not send an CONECT-Packet")
		err = errors.New("Client not send an CONECT-Packet")
		return
	}

	switch recvPacket := eventPacket.(type) {

	case *packets.SubscribePacket:
		err = broker.onPacketSubscribe(eventClient, recvPacket)

	case *packets.UnsubscribePacket:
		err = broker.onPacketUnSubscribe(eventClient, recvPacket)

	case *packets.PingreqPacket:
		err = broker.onPacketPinReq(eventClient, recvPacket)

	case *packets.PublishPacket:
		err = broker.onPacketPublish(eventClient, recvPacket)

	case *packets.PubackPacket:
		err = broker.onPacketPubACK(eventClient, recvPacket)

	case *packets.PubrecPacket:
		err = broker.onPacketPubRec(eventClient, recvPacket)

	case *packets.PubrelPacket:
		err = broker.onPacketPubRel(eventClient, recvPacket)

	case *packets.PubcompPacket:
		err = broker.onPacketPubcomp(eventClient, recvPacket)

	}

	return
}
