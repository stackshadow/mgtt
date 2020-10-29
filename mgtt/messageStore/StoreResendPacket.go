package messagestore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

type StoreResendPacketOptions struct {
	ClientID string
	OriginID uint16
	ResendAt time.Time
	Packet   *packets.PublishPacket
}

type packetInfo struct {
	ClientID   string    `json:"c,omitempty"`
	OriginID   uint16    `json:"o,omitempty"`
	ResendAt   time.Time `json:"t"`
	PacketData []byte    `json:"p"` // this should not be set from outside, it will be overwritten
}

// StoreResendPacket will store an `packet` to resend it after `delay`
//
// If IDStart is already used, a new free id will returned in IDUsed
//
// if IDStart is free, IDUsed = IDStart
func (store *Store) StoreResendPacket(bucket string, lastUsedID *uint16, options *StoreResendPacketOptions) (err error) {

	var newPacketInfo packetInfo

	// store already known infos
	newPacketInfo.ClientID = options.ClientID
	newPacketInfo.OriginID = options.Packet.MessageID
	newPacketInfo.ResendAt = options.ResendAt

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
			// create id from uint16
			binary.LittleEndian.PutUint16(newIDBytes, *lastUsedID)

			existingPacket := b.Get(newIDBytes)
			if existingPacket == nil {
				break
			}

			*lastUsedID++
		}

		// okay, found a free ID, we store it to the packet
		options.Packet.MessageID = *lastUsedID

		// convert packet to bytes
		writer := bytes.NewBuffer([]byte{})
		options.Packet.Write(writer)
		newPacketInfo.PacketData = writer.Bytes()

		// create json-byte-array
		var payload []byte
		payload, err = json.Marshal(newPacketInfo)

		// save it to DB
		err = b.Put(newIDBytes, payload)

		return err
	})

	if err != nil {
		log.Error().
			Uint16("mid", options.Packet.MessageID).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store packet")
	} else {
		log.Debug().
			Uint16("mid", options.Packet.MessageID).
			Str("bucket", bucket).
			Msg("Packet stored")
	}

	return
}
