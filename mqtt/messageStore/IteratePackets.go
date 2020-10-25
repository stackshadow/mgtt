package messagestore

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// IteratePackets will iterate published-packet in a bucket
func (store *Store) IteratePackets(bucket string, iterate func(packet *packets.PublishPacket)) (err error) {

	err = store.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// payload
			publishPacketData := bytes.NewBuffer(v)

			// load the packet
			publishPacketGeneric, _ := packets.ReadPacket(publishPacketData)

			// convert it
			publishPacket := publishPacketGeneric.(*packets.PublishPacket)

			// call iterate-function
			iterate(publishPacket)
		}

		return nil
	})

	return
}
