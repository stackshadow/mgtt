package messagestore

import (
	"github.com/boltdb/bolt"
)

type Store struct {
	db *bolt.DB
}

// Open will open the DB and create needed buckets
func Open(dbPath string) (store *Store, err error) {

	store = &Store{}

	store.db, err = bolt.Open(dbPath, 0600, nil)

	return
}

// Close an opened DB
func (store *Store) Close() {
	store.db.Close()
}
