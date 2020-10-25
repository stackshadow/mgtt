package messagestore

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// StorePacketWithTopic will store an published-packet in a bucket
//
// Packet can be overwritten !
func (store *Store) StorePacketWithTopic(bucket string, topic string, packet *packets.PublishPacket) (err error) {

	// payload
	writer := bytes.NewBuffer([]byte{})
	packet.Write(writer)
	payload := writer.Bytes()

	// save it to the db
	err = store.db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if b == nil {
			return nil
		}

		// save
		err = b.Put([]byte(topic), payload)

		return err
	})

	if err != nil {
		log.Error().
			Str("topic", topic).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store payload to retained-store")
	} else {
		log.Debug().
			Str("topic", topic).
			Str("bucket", bucket).
			Msg("Stored payload to retained-store")
	}

	return
}
