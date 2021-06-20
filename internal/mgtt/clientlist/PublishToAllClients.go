package clientlist

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

// PublishToAllClients send an packet to all clients
func PublishToAllClients(packet *packets.PublishPacket, skipClientID string, once bool) (published bool, subscribed bool, err error) {

	// mutex
	listMutex.Lock()
	defer listMutex.Unlock()

	for _, client := range list {

		clientID := client.ID()
		userName := client.Username()

		if clientID == skipClientID {
			continue
		}

		if plugin.CallOnSendToSubscriberRequest(clientID, userName, packet.TopicName) == true {

			publishOk, subscriptionOK, publishErr := client.Publish(packet)
			if once == true {
				if publishOk == true {
					return true, true, nil
				}
			}

			published = published || publishOk
			subscribed = subscribed || subscriptionOK

			if publishErr != nil && err == nil {
				err = publishErr
			}

		}

	}

	return
}
