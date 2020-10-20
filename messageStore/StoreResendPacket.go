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

// PacketResendInfo hold information about a package that should be resended
type StoreResendPacketOption struct {
	BrokerMessageID uint16
	ClientID        string
	ResendAt        time.Time
	Packet          *packets.PublishPacket
}

type packetInfo struct {
	ClientID   string    `json:"s"`
	ResendAt   time.Time `json:"t"`
	PacketData []byte    `json:"p"` // this should not be set from outside, it will be overwritten
}

// StoreResendPacket will store an `packet` to resend it after `delay`
//
// If IDStart is already used, a new free id will returned in IDUsed
//
// if IDStart is free, IDUsed = IDStart
func (store *Store) StoreResendPacket(bucket string, option *StoreResendPacketOption) (err error) {

	// convert packet to bytes
	writer := bytes.NewBuffer([]byte{})
	option.Packet.Write(writer)

	// packetinfo
	newPacketInfo := packetInfo{
		ClientID:   option.ClientID,
		ResendAt:   option.ResendAt,
		PacketData: writer.Bytes(),
	}

	// some ne
	newIDBytes := make([]byte, 2)

	// save it to the db
	err = store.db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if b == nil {
			return nil
		}

		// try to find a free ID
		for {
			// create id from uint16
			binary.LittleEndian.PutUint16(newIDBytes, option.BrokerMessageID)

			existingPacket := b.Get(newIDBytes)
			if existingPacket == nil {
				break
			}

			option.BrokerMessageID++
		}

		// create json-byte-array
		var payload []byte
		payload, err = json.Marshal(newPacketInfo)
		err = b.Put(newIDBytes, payload)

		return err
	})

	if err != nil {
		log.Error().
			Uint16("packet-mid", option.Packet.MessageID).
			Uint16("broker-mid", option.BrokerMessageID).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store packet")
	} else {
		log.Debug().
			Uint16("packet-mid", option.Packet.MessageID).
			Uint16("broker-mid", option.BrokerMessageID).
			Str("bucket", bucket).
			Msg("Packet stored")
	}

	return
}
