package client

import "github.com/eclipse/paho.mqtt.golang/packets"

// SendPingresp will send an PINGRESP-Package
func (evt *Event) SendPingresp() {

	// convert
	_, ok := evt.Packet.(*packets.PingreqPacket)
	if ok == false {
		return
	}

	pingResp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
	pingResp.Write(evt.Client.connection)

	return
}
