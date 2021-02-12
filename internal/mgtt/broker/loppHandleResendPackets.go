package broker

import (
	"errors"
	"net"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) loopHandleResendPackets() (err error) {

	var retryClient *client.MgttClient = &client.MgttClient{}

	netserver, _ := net.Pipe()
	retryClient.Init(netserver, 0)
	retryClient.IDSet("resend")
	retryClient.Connected = true

	triggerTime := time.Now()
	triggerTime = triggerTime.Add(time.Minute * 1)

loop:
	for {

		// check if we need to resend messages that are not replyed with PUBACK
		log.Debug().Msg("Check if packets need to be resended")

		select {
		case <-broker.loopHandleResendPacketsExit:
			log.Debug().Msg("STOP go-routine loopHandleResendPackets()")
			break loop
		default:

			if time.Now().After(triggerTime) == false {
				time.Sleep(time.Minute * 1)
				continue
			}
			triggerTime = time.Now()
			triggerTime = triggerTime.Add(time.Minute * 1)

			persistance.PacketIterate("qos", func(info persistance.PacketInfo) {
				// check if time is up
				if time.Now().After(info.RetryAt) == true {

					// set the time to the future
					info.RetryAt = time.Now().Add(time.Minute)
					persistance.PacketStore("qos", &info)

					if info.PubRec == false {

						// create a packet
						pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
						pubPacket.MessageID = info.MessageID
						pubPacket.Retain = false
						pubPacket.Dup = info.OriginClientID != ""
						pubPacket.TopicName = info.Topic
						pubPacket.Payload = info.Payload
						pubPacket.Qos = info.Qos

						// publish packet to all subscribers
						switch info.Qos {
						case 1:
							_, _, err = broker.PublishPacket(pubPacket, false)
						case 2:
							_, _, err = broker.PublishPacket(pubPacket, true)
						default:
							err = errors.New("This QOS-Level is not supported")
						}

					} else {
						clientlist.PubrelToClient(info.TargetClientID, info.MessageID)
					}
				}

			})

		}

		time.Sleep(time.Minute * 1)
	}

	return
}
