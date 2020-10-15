package messagestore

import (
	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

// Open will open the DB and create needed buckets
func Open() (store *Store, err error) {

	store = &Store{}

	store.db, err = bolt.Open("messages.db", 0600, nil)

	store.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("retainedTopics"))
		return err
	})

	return
}

// Close an opened DB
func (store *Store) Close() {
	store.db.Close()
}
