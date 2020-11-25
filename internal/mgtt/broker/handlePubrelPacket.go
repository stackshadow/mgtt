package broker

import (
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/client"
)

func (broker *Broker) handlePubrelPacket(connectedClient *client.MgttClient, packet *packets.PubrelPacket) (err error) {

	// find the package
	storedInfo, err := broker.retainedMessages.FindPacket("unreleased", connectedClient.ID(), packet.MessageID)
	if err != nil {
		return err
	}

	// and store it to "resend"
	broker.lastIDLock.Lock()

	storedInfo.ResendAt = time.Now().Add(time.Minute * 1)
	storedInfo.Packet.Dup = true
	err = broker.retainedMessages.StoreResendPacket("resend", &broker.lastID, storedInfo)
	if err != nil {
		return err
	}

	broker.lastIDLock.Unlock()

	// remove it from unreleased
	broker.retainedMessages.DeletePacketWithID("unreleased", packet.MessageID)

	// [MQTT-3.3.1-9]
	// MUST set the RETAIN flag to 0 when a PUBLISH Packet is sent to a Client
	// because it matches an established subscription
	storedInfo.Packet.Retain = false

	// WE publish exact once !
	_, err = broker.PublishPacket(storedInfo.Packet, true)

	// we dont do this
	// we wait for pubrec then we notify
	//if published == true {
	//	event.client.SendPubcomp(pubrelPacket.MessageID)
	//}

	return
}
