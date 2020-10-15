package messagestore

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// IterateRetainedTopics will call the iterate-function for every retained topic
func (store *Store) IterateRetainedTopics(iterate func(packet *packets.PublishPacket)) (err error) {

	store.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("retainedTopics"))

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
