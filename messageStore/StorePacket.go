package messagestore

import (
	"bytes"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// storePacket will store an published-packet in a bucket
//
// If the id already exist an error will be returned
func (store *Store) storePacket(bucket string, key []byte, packet *packets.PublishPacket) (err error) {

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

		// only save if not exist
		existingPacket := b.Get(key)
		if existingPacket != nil {
			err = errors.New("Packet already exist in bucket, can not be overwritten")
		} else {
			err = b.Put(key, payload)
		}

		return err
	})

	return
}
