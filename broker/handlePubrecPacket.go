package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) handlePubrecPacket(event *Event) (err error) {

	// check package
	pubrecPacket, ok := event.packet.(*packets.PubrecPacket)
	if ok == false {
		err = errors.New("Package is not packets.PublishPacket")
		return
	}

	// get infos
	var storedInfo *messagestore.StoreResendPacketOption
	storedInfo, err = broker.retainedMessages.GetResendPacket("resend", pubrecPacket.MessageID)
	if err == nil {
		broker.pubrec[storedInfo.BrokerMessageID] = storedInfo
		broker.retainedMessages.DeletePacketWithID("resend", storedInfo.BrokerMessageID)
	}

	return
}
