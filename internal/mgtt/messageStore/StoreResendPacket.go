package messagestore

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// PacketInfo represent the struct which will be stored in the bolt-db
type PacketInfo struct {
	ResendAt time.Time `json:"r"` //

	// tis comes from the packet
	Topic     string `json:"t,omitempty"`
	MessageID uint16 `json:"i,omitempty"`
	Qos       byte   `json:"q,omitempty"`
	Payload   []byte `json:"d,omitempty"`
}

// StoreResendPacket will store an `packet` to resend it after `delay`
//
// If IDStart is already used, a new free id will returned in IDUsed
//
// if IDStart is free, IDUsed = IDStart
func (store *Store) StoreResendPacket(bucket string, info *PacketInfo) (err error) {

	// save it to the db
	err = store.db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if b == nil {
			return nil
		}

		// try to find a free ID
		newIDBytes := make([]byte, 2)
		for {
			// convert to le-uint16
			binary.LittleEndian.PutUint16(newIDBytes, info.MessageID)

			existingPacket := b.Get(newIDBytes)
			if existingPacket == nil {
				break
			}

			info.MessageID++
		}

		// create json-byte-array
		var payload []byte
		payload, err = json.Marshal(info)

		// save it to DB
		err = b.Put(newIDBytes, payload)

		return err
	})

	if err != nil {
		log.Error().
			Uint16("mid", info.MessageID).
			Str("topic", info.Topic).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store packet")
	} else {
		log.Debug().
			Uint16("mid", info.MessageID).
			Str("topic", info.Topic).
			Str("bucket", bucket).
			Msg("Packet stored")
	}

	return
}
