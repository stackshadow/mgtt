package messagestore

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// IterateResendPackets will iterate packages that are stored with StoreResendPacket()
func (store *Store) IterateResendPackets(bucket string, iterate func(storedInfo *PacketInfo)) (err error) {

	err = store.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// parse json
			var info PacketInfo
			err = json.Unmarshal(v, &info)
			if err != nil {
				return err
			}

			// call iterate-function
			iterate(&info)
		}

		return nil
	})

	return
}
