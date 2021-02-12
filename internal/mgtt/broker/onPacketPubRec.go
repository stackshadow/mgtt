package broker

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/client"
	"gitlab.com/mgtt/internal/mgtt/persistance"
)

func (broker *Broker) onPacketPubRec(connectedClient *client.MgttClient, packet *packets.PubrecPacket) (err error) {

	log.Debug().
		Str("cid", connectedClient.ID()).
		Uint16("pid", packet.MessageID).
		Msg("Get pubrec, reply with pubrel")

	// get packet-info
	var packetInfo persistance.PacketInfo
	var packetInfoExist bool = false
	packetInfoExist, packetInfo, _ = persistance.PacketExist("qos", persistance.PacketFindOpts{
		MessageID: &packet.MessageID,
	})

	// store that we get pubrec
	if packetInfoExist == true {
		packetInfo.TargetClientID = connectedClient.ID()
		packetInfo.PubRec = true
		persistance.PacketStore("qos", &packetInfo)
	}

	// request pubrel
	connectedClient.SendPubrel(packet.MessageID)

	return
}
