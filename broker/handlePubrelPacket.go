package broker

import (
	"errors"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

func (broker *Broker) handlePubrelPacket(event *Event) (err error) {

	// check package
	pubrelPacket, ok := event.packet.(*packets.PubrelPacket)
	if ok == false {
		err = errors.New("Package is not packets.PubrelPacket")
		return
	}

	// find the package
	storedInfo, err := broker.retainedMessages.FindPacket("unreleased", event.client.ID(), pubrelPacket.MessageID)
	if err != nil {
		return err
	}
	messageIDInUnreleased := storedInfo.BrokerMessageID

	// and store it to "resend"
	broker.lastIDLock.Lock()
	storedInfo.ResendAt = time.Now().Add(time.Minute * 1)
	storedInfo.Packet.Dup = true
	storedInfo.BrokerMessageID = broker.lastID + 1

	err = broker.retainedMessages.StoreResendPacket("resend", storedInfo)
	if err != nil {
		return err
	}
	broker.lastID = storedInfo.BrokerMessageID
	broker.lastIDLock.Unlock()

	// remove it from unreleased
	broker.retainedMessages.DeletePacketWithID("unreleased", messageIDInUnreleased)

	// because we stored the original message with the original messageID, we can now manipulate it
	storedInfo.Packet.MessageID = storedInfo.BrokerMessageID

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
	// because it matches an established subscription
	storedInfo.Packet.Retain = false

	// WE publish exact once !
	var published bool
	for _, client := range broker.clients {
		published, err = client.Publish(storedInfo.Packet)
		if published == true {
			break
		}
	}

	// we dont do this
	// we wait for pubrec then we notify
	//if published == true {
	//	event.client.SendPubcomp(pubrelPacket.MessageID)
	//}

	return
}
