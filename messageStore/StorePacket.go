package messagestore

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// StorePacket will store an published-packet in a bucket
func (store *Store) StorePacket(bucket string, id string, packet *packets.PublishPacket) (err error) {

	// payload
	writer := bytes.NewBuffer([]byte{})
	packet.Write(writer)
	payload := writer.Bytes()

	// save it to the db
	err = store.db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		err = b.Put([]byte(id), payload)
		return err
	})

	if err != nil {
		log.Debug().
			Str("id", id).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store payload to retained-store")
	} else {
		log.Debug().
			Str("id", id).
			Str("bucket", bucket).
			Msg("Stored payload to retained-store")
	}

	return
}
