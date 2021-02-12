package persistance

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

// PacketIterate will iterate packets
func PacketIterate(bucket string, iterate func(info PacketInfo)) (err error) {

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// payload
			var packetInfo PacketInfo

			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				// call iterate-function
				iterate(packetInfo)
			}
		}

		return nil
	})

	return
}
