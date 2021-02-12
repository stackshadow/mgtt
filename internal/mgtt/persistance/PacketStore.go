package persistance

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

// PacketStore will store an packetInfo to persistance store
//
// - need Open() before
//
// - if info.MessageID == 0 then a new messageID will be used otherwise the message gets overwritten
//
// - if lastID exist in the store, it will be set to the next available number
//
// - This function is thread-save over mutex
func PacketStore(bucket string, info *PacketInfo) (err error) {

	packetStoreMutex.Lock()
	defer packetStoreMutex.Unlock()

	// save it to the db
	err = db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(bucket))
		if b == nil {
			return nil
		}

		// try to find a free ID
		newIDBytes := make([]byte, 2)
		if info.MessageID == 0 {
			for {
				// convert to le-uint16
				binary.LittleEndian.PutUint16(newIDBytes, packetStoreLastID)

				existingPacket := b.Get(newIDBytes)
				if existingPacket == nil {
					break
				}

				packetStoreLastID++
			}

			// store new id
			info.MessageID = packetStoreLastID
		} else {
			binary.LittleEndian.PutUint16(newIDBytes, info.MessageID)
		}

		// create json-byte-array
		if payload, err := json.Marshal(info); err == nil {
			// save it to DB
			err = b.Put(newIDBytes, payload)
			info.dump(bucket, "Stored")
		}

		return err
	})

	packetStoreDump(bucket)
	return
}
