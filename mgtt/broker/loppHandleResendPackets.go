package broker

import (
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	messagestore "gitlab.com/mgtt/messageStore"
)

func (broker *Broker) loopHandleResendPackets() {

	netserver, _ := net.Pipe()
	retryClient := client.New(netserver, 0)
	retryClient.IDSet("resend")
	retryClient.Connected = true

	go func() {
		for {

			// wait a bit
			time.Sleep(time.Minute * 1)

			// check if we need to resend messages that are not replyed with PUBACK
			log.Debug().Msg("Check if packets need to be resended")
			broker.retainedMessages.IterateResendPackets("resend", func(storedInfo *messagestore.StoreResendPacketOptions) {

				// check if time is up
				if time.Now().After(storedInfo.ResendAt) == true {

					storedInfo.Packet.Dup = true     // this is an duplicate packet
					storedInfo.Packet.Retain = false // resend, not retain ;)

					log.Debug().
						Uint16("packet-mid", storedInfo.Packet.MessageID).
						Str("topic", storedInfo.Packet.TopicName).
						Msg("Resend packet")

					normalClose, err := broker.loopHandleBrokerPacket(retryClient, storedInfo.Packet)
					if err != nil || normalClose {
						return
					}

					// a small delay to not flood our clients
					time.Sleep(time.Millisecond * 500)

				}

			})

		}
	}()
}
