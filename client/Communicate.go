package client

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Communicate reads the packets from the io.Writer and push them to the chan
//
// return err if an error occured
func (c *MgttClient) Communicate(clientEvents chan *Event) (err error) {

	for {
		var commonpacket packets.ControlPacket
		commonpacket, err = packets.ReadPacket(c.connection)
		if err != nil {
			break
		}
		clientEvents <- &Event{
			Client: c,
			Packet: commonpacket,
		}

	}

	return
}
