package clientlist

import (
	"net"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Client  represents a client
type Client interface {
	Init(connection net.Conn, timeout time.Duration)
	ResetTimeout()
	Close() (err error)
	RemoteAddr() (remoteAddr string) // return the remoteAddr as string

	// ID's
	IDSet(id string)
	ID() string

	UsernameSet(username string)
	Username() string

	// Last-will
	LastWillSet(packet *packets.PublishPacket)
	LastWillGet() (packet *packets.PublishPacket)

	// packet handling
	PacketRead() (packet packets.ControlPacket, err error) // reads a single packet from the connection
	Publish(packet *packets.PublishPacket) (published bool, subscribed bool, err error)

	// subscriptions
	SubScriptionAdd(topic string)
	SubScriptionsAdd(topics []string)

	//
	SendConnack(ReturnCode byte, SessionPresent bool) (err error)
	SendConnect(username, password, clientid string)
	SendPingreq()
	SendPingresp()
	SendPuback(MessageID uint16) (err error)
	SendPubcomp(MessageID uint16) (err error)
	SendPubrec(MessageID uint16) (err error)
	SendPubrel(MessageID uint16) (err error)
	SendSuback(MessageID uint16, ReturnCodes []byte) (err error)
}
