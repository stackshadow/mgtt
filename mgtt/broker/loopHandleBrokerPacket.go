package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
)

// Communicate will handle incoming messages
//
// - this is a BLOCKING function
func (broker *Broker) loopHandleBrokerPacket(eventClient *client.MgttClient, eventPacket packets.ControlPacket) (normalClose bool, err error) {

	log.Debug().
		Uint16("mid", eventPacket.Details().MessageID).
		Uint8("Qos", eventPacket.Details().Qos).
		Str("packet", eventPacket.String()).
		Msg("Handle packet")

	switch recvPacket := eventPacket.(type) {
	case *packets.ConnectPacket:
		err = broker.handleConnectPacket(eventClient, recvPacket)
		return
	case *packets.DisconnectPacket:
		err = broker.handleDisConnectPacket(eventClient)
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
		err = broker.handleSubscribePacket(eventClient, recvPacket)

	case *packets.PingreqPacket:
		err = broker.handlePingreqPacket(eventClient, recvPacket)

	case *packets.PublishPacket:
		err = broker.handlePublishPacket(eventClient, recvPacket)

	case *packets.PubackPacket:
		err = broker.handlePubacPacket(eventClient, recvPacket)

	case *packets.PubrecPacket:
		err = broker.handlePubrecPacket(eventClient, recvPacket)

	case *packets.PubrelPacket:
		err = broker.handlePubrelPacket(eventClient, recvPacket)

	case *packets.PubcompPacket:
		err = broker.handlePubcompPacket(eventClient, recvPacket)

	}

	return
}
