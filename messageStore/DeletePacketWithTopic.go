package messagestore

import (
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// DeletePacketWithTopic will delete an published-packet in a bucket if it exist
//
// this function only return an error Returns an error if the bucket was
// created from a read-only transaction
func (store *Store) DeletePacketWithTopic(bucket string, topic string) (err error) {

	err = store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("retainedTopics"))
		err = b.Delete([]byte(topic))

		if err == nil {

			log.Debug().
				Str("topic", topic).
				Msg("Delete payload from retained-store")
		}

		return err
	})

	return
}
