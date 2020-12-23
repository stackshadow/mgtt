package clientlist

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

// PublishToAllClients send an packet to all clients
func PublishToAllClients(packet *packets.PublishPacket, once bool) (published bool, err error) {

	for _, client := range list {

		clientID := client.ID()
		userName := client.Username()
		if plugin.CallOnSendToSubscriberRequest(clientID, userName, packet.TopicName) == true {

			publishOk, publishErr := client.Publish(packet)
			if once == true {
				if publishOk == true {
					return true, nil
				}
			}

			published = published || publishOk
			if publishErr != nil && err == nil {
				err = publishErr
			}

		}

	}

	return
}
