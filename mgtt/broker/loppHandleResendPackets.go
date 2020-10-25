package broker

import (
	"time"

	"github.com/rs/zerolog/log"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) loopHandleResendPackets() {

	for {

		// wait a bit
		time.Sleep(time.Minute * 1)

		// check if we need to resend messages that are not replyed with PUBACK
		log.Debug().Msg("Check if packets need to be resended")
		broker.retainedMessages.IterateResendPackets("resend", func(storedInfo *messagestore.StoreResendPacketOption) {

			// check if time is up
			if time.Now().After(storedInfo.ResendAt) == true {

				storedInfo.Packet.Dup = true                             // this is an duplicate packet
				storedInfo.Packet.Retain = false                         // resend, not retain ;)
				storedInfo.Packet.MessageID = storedInfo.BrokerMessageID // we use our message ID

				log.Debug().
					Uint16("packet-mid", storedInfo.Packet.MessageID).
					Uint16("broker-mid", storedInfo.BrokerMessageID).
					Str("topic", storedInfo.Packet.TopicName).
					Msg("Resend packet")

				// new event
				newEvent := Event{
					client: nil,
					packet: storedInfo.Packet,
				}

				broker.clientEvents <- &newEvent
			}

		})

	}

}
