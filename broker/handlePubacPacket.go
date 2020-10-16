package broker

import (
	"errors"
	"fmt"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/client"
)

func (broker *Broker) handlePubacPacket(event *client.Event) (err error) {

	// check package
	packet, ok := event.Packet.(*packets.PubackPacket)
	if ok == false {
		err = errors.New("Package is not packets.PubackPacket")
		return
	}

	broker.retainedMessages.DeletePacketWithTopic("resend", fmt.Sprintf("%d", packet.MessageID))

	return
}
