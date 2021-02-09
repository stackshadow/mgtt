package persistance

import (
	"time"

	"github.com/rs/zerolog/log"
)

// PacketInfo represent the struct which will be stored in the bolt-db
type PacketInfo struct {
	// retry-Timestamp
	RetryAt time.Time `json:"r"`

	// Original infos
	OriginClientID  string `json:"oc"`
	TargetClientID  string `json:"tc,omitempty"` // This is for retry of PUBREL
	OriginMessageID uint16 `json:"om,omitempty"`

	// tis comes from the packet
	MessageID uint16 `json:"i,omitempty"` // If 0, a new free messageID will be assigned
	Topic     string `json:"t,omitempty"`
	Qos       byte   `json:"q,omitempty"`
	Payload   []byte `json:"d,omitempty"`

	// states
	PubRec bool `json:"p"` // pubrec received, stop retry publish
}

func (info PacketInfo) dump(bucket, message string) {

	logger := log.Logger.With().CallerWithSkipFrameCount(3).Logger()

	logger.Debug().
		Str("bucket", bucket).
		Str("ocid", info.OriginClientID).
		Uint16("omid", info.OriginMessageID).
		Uint16("mid", info.MessageID).
		Str("topic", info.Topic).
		Msg(message)
}
