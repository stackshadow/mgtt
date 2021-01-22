package persistance

import "time"

// PacketInfo represent the struct which will be stored in the bolt-db
type PacketInfo struct {
	OriginClientID string    `json:"o"`
	ResendAt       time.Time `json:"r"` //

	// tis comes from the packet
	Topic     string `json:"t,omitempty"`
	MessageID uint16 `json:"i,omitempty"`
	Qos       byte   `json:"q,omitempty"`
	Payload   []byte `json:"d,omitempty"`
}
