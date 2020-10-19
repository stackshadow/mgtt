package messagestore

import (
	"bytes"
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// StorePacketWithID will store an published-packet in a bucket
//
// If 'id' exists:
//
// - and the packet behind this id is NIL, your packet is stored and 'newID' = 'id'
//
// - and an packet exist behind this, a new free ID will be searched
//
// If 'packet' is NIL, a new free ID will be searched
//
// This "feature" make it possible to reserve an ID and overwrite it later
//
// so, to reserve an ID, use 'packet' = NIL
func (store *Store) StorePacketWithID(bucket string, id uint16, packet *packets.PublishPacket) (newID uint16, err error) {

	// payload
	var payload []byte
	if packet != nil {
		writer := bytes.NewBuffer([]byte{})
		packet.Write(writer)
		payload = writer.Bytes()
	} else {
		payload = []byte{0}
	}

	newID = id
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
			binary.LittleEndian.PutUint16(newIDBytes, uint16(newID))

			existingPacket := b.Get(newIDBytes)
			if existingPacket == nil {
				break
			}

			// if key exist with a 1-Byte-Value, its an reserved ID and we can write to it :)
			// but only if we not search for a new ID
			if existingPacket != nil {
				if len(existingPacket) == 1 && id == newID {
					break
				}
			}

			newID++
		}

		err = b.Put(newIDBytes, payload)

		return err
	})
	if packet != nil {
		if err != nil {
			log.Error().
				Uint16("id", newID).
				Str("bucket", bucket).
				Err(err).
				Msg("Can not store payload to retained-store")
		} else {
			log.Debug().
				Uint16("id", newID).
				Str("bucket", bucket).
				Msg("Stored payload to retained-store")
		}
	} else {
		if err != nil {
			log.Error().
				Uint16("id", newID).
				Str("bucket", bucket).
				Err(err).
				Msg("Can not reserve id")
		} else {
			log.Debug().
				Uint16("id", newID).
				Str("bucket", bucket).
				Msg("Reserve id")
		}
	}

	return
}
