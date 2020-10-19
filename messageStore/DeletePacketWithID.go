package messagestore

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// DeletePacketWithTopic will delete an published-packet in a bucket if it exist
//
// this function only return an error Returns an error if the bucket was
// created from a read-only transaction
func (store *Store) DeletePacketWithID(bucket string, id uint16) (err error) {

	newIDBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(newIDBytes, id)

	err = store.db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		err = b.Delete(newIDBytes)

		if err == nil {

			log.Debug().
				Uint16("id", id).
				Str("bucket", bucket).
				Msg("Delete payload from retained-store")
		}

		return err
	})

	return
}
