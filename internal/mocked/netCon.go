package mocked

import (
	"net"
	"time"
)

type Con struct {
	packetSendLoopExit chan byte
}

func ConNew() (c Con) {

	return Con{
		packetSendLoopExit: make(chan byte),
	}

}

func (c Con) Read(b []byte) (n int, err error) {

	bytes := <-c.packetSendLoopExit
	b[0] = bytes
	return 1, nil
}
func (c Con) Write(b []byte) (n int, err error) {

	for _, singleByte := range b {
		c.packetSendLoopExit <- singleByte
	}

	return len(b), nil
}

func (c Con) Close() error {
	return nil
}

func (c Con) LocalAddr() (add net.Addr) {
	return
}

func (c Con) RemoteAddr() (add net.Addr) {
	return
}

func (c Con) SetDeadline(t time.Time) error {
	return nil
}

func (c Con) SetReadDeadline(t time.Time) error {
	return nil
}

func (c Con) SetWriteDeadline(t time.Time) error {
	return nil
}
