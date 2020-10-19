package messagestore

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// StorePacketWithTopic will store an published-packet in a bucket
//
// If the id already exist an error will be returned
func (store *Store) StorePacketWithTopic(bucket string, topic string, packet *packets.PublishPacket) (err error) {

	err = store.storePacket(bucket, []byte(topic), packet)

	if err != nil {
		log.Error().
			Str("topic", topic).
			Str("bucket", bucket).
			Err(err).
			Msg("Can not store payload to retained-store")
	} else {
		log.Debug().
			Str("topic", topic).
			Str("bucket", bucket).
			Msg("Stored payload to retained-store")
	}

	return
}
