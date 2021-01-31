package persistance

import (
	"encoding/json"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// PacketFindOpts represents options for variouse search-functions
type PacketFindOpts struct {
	OriginClientID  *string
	OriginMessageID *uint16
	MessageID       *uint16
	Topic           *string
}

var packetStoreLastID uint16 = 1
var packetStoreMutex sync.Mutex

func packetStoreDump(bucket string) {
	return

	log.Debug().Str("bucket", bucket).Msg("###################### DUMP of packet-store ######################")

	PacketIterate(bucket, func(info PacketInfo, publishPacket *packets.PublishPacket) {
		info.dump(bucket, "")
	})

	log.Debug().Str("bucket", bucket).Msg("###################### DUMP end ######################")
}

func packetGet(bucket string, opts PacketFindOpts) (found bool, key []byte, packetInfo PacketInfo, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				found = true

				if opts.OriginClientID != nil {
					found = packetInfo.OriginClientID == *opts.OriginClientID && found
				}

				if opts.OriginMessageID != nil {
					found = packetInfo.OriginMessageID == *opts.OriginMessageID && found
				}

				if opts.MessageID != nil {
					found = packetInfo.MessageID == *opts.MessageID && found
				}

				if opts.Topic != nil {
					found = packetInfo.Topic == *opts.Topic && found
				}

				if found == true {
					key = k
					packetInfo.dump(bucket, "Found packet")
					break
				}

			}

		}

		return err
	})

	return
}
