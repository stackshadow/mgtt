package persistance

import "github.com/boltdb/bolt"

// Open will open the DB
func Open(dbFilePath string) (err error) {
	db, err = bolt.Open(dbFilePath, 0600, nil)
	return
}
