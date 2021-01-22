package persistance

import "github.com/boltdb/bolt"

// Open will open the DB
func Open(dbPath string) (err error) {
	db, err = bolt.Open(dbPath, 0600, nil)
	return
}
