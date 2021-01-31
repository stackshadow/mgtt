package persistance

import "github.com/boltdb/bolt"

// PacketDelete will delete the first packet it found
//
// - If more than one parameter is not NIL all must match
//
// - This function is thread-save over mutex
func PacketDelete(bucket string, opts PacketFindOpts) (err error) {

	packetStoreMutex.Lock()
	defer packetStoreMutex.Unlock()

	var found bool
	var foundKey []byte
	var foundPacketInfo PacketInfo

	for {
		found, foundKey, foundPacketInfo, err = packetGet(bucket, opts)

		if found && err == nil {
			err = db.Update(func(tx *bolt.Tx) error {
				// Assume bucket exists and has keys
				b := tx.Bucket([]byte(bucket))
				if b == nil {
					return nil
				}

				if err := b.Delete(foundKey); err != nil {
					foundPacketInfo.dump(bucket, "Deleted")
				}
				return err
			})
		} else {
			break
		}
	}

	packetStoreDump(bucket)
	return
}
