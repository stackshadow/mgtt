package messagestore

import (
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// DeleteRetainedIfExist will delete an topic if it exist
func (store *Store) DeleteRetainedIfExist(topic string) (err error) {

	store.db.Update(func(tx *bolt.Tx) error {
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
