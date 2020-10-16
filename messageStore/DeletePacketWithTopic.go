package messagestore

import (
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// DeletePacketWithTopic will delete an published-packet in a bucket if it exist
//
// this function only return an error Returns an error if the bucket was
// created from a read-only transaction
func (store *Store) DeletePacketWithTopic(bucket string, id string) (err error) {

	err = store.db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		err = b.Delete([]byte(id))

		if err == nil {

			log.Debug().
				Str("id", id).
				Str("bucket", bucket).
				Msg("Delete payload from retained-store")
		}

		return err
	})

	return
}
