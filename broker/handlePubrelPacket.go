package broker

import (
	"errors"

	"github.com/eclipse/paho.mqtt.golang/packets"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) handlePubrelPacket(event *Event) (err error) {

	// check package
	pubrelPacket, ok := event.packet.(*packets.PubrelPacket)
	if ok == false {
		err = errors.New("Package is not packets.PubrelPacket")
		return
	}

	// get infos
	var publishPacket *messagestore.StoreResendPacketOption
	var brokerMessageID uint16
	for brokerMessageID, publishPacket = range broker.pubrec {
		if event.client.ID() == publishPacket.ClientID {
			if pubrelPacket.MessageID == publishPacket.Packet.MessageID {
				broker.retainedMessages.DeletePacketWithID("resend", brokerMessageID)
				event.client.SendPubcomp(pubrelPacket.MessageID)
			}
		}
	}

	return
}
