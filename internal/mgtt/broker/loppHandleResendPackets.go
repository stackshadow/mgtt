package broker

import (
	"net"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	messagestore "gitlab.com/mgtt/internal/mgtt/messageStore"
)

func (broker *Broker) loopHandleResendPackets() {

	netserver, _ := net.Pipe()
	retryClient := client.New(netserver, 0)
	retryClient.IDSet("resend")
	retryClient.Connected = true

	triggerTime := time.Now()
	triggerTime = triggerTime.Add(time.Minute * 1)

loop:
	for {

		select {
		case <-broker.loopHandleResendPacketsExit:
			log.Debug().Msg("STOP go-routine loopHandleResendPackets()")
			break loop
		default:

			if time.Now().After(triggerTime) == false {
				time.Sleep(time.Second * 1)
				continue
			}
			triggerTime = time.Now()
			triggerTime = triggerTime.Add(time.Minute * 1)

			// check if we need to resend messages that are not replyed with PUBACK
			log.Debug().Msg("Check if packets need to be resended")
			broker.retainedMessages.IterateResendPackets("resend", func(storedInfo *messagestore.PacketInfo) {

				// check if time is up
				if time.Now().After(storedInfo.ResendAt) == true {

					// we create a new publish packet
					pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
					pubPacket.MessageID = storedInfo.MessageID
					pubPacket.Retain = false
					pubPacket.Dup = true // this is an duplicate packet
					pubPacket.TopicName = storedInfo.Topic
					pubPacket.Payload = storedInfo.Payload
					pubPacket.Qos = storedInfo.Qos

					log.Debug().
						Uint16("packet-mid", pubPacket.MessageID).
						Str("topic", pubPacket.TopicName).
						Msg("Resend packet")

					err := broker.handlePublishPacket(retryClient, pubPacket)
					if err != nil {
						return
					}

					// a small delay to not flood our clients
					time.Sleep(time.Millisecond * 500)

				}

			})

		}
	}

}
