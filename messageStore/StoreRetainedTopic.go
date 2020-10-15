package messagestore

import (
	"bytes"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// StoreRetainedTopic will save payload to the "retainedTopics"-Bucket to topic
func (store *Store) StoreRetainedTopic(packet *packets.PublishPacket) (err error) {

	// topic
	topic := packet.TopicName

	// payload
	writer := bytes.NewBuffer([]byte{})
	packet.Write(writer)
	payload := writer.Bytes()

	// save it to the db
	err = store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("retainedTopics"))
		err = b.Put([]byte(topic), payload)
		return err
	})

	if err != nil {
		log.Debug().
			Str("topic", topic).
			Err(err).
			Msg("Can not store payload to retained-store")
	} else {
		log.Debug().
			Str("topic", topic).
			Msg("Stored payload to retained-store")
	}

	return
}
