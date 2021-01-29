package persistance

import (
	"encoding/binary"
	"encoding/json"
	"sync"

	"time"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/rs/zerolog/log"
)

// PacketInfo represent the struct which will be stored in the bolt-db
type PacketInfo struct {
	OriginClientID string    `json:"o"`
	ResendAt       time.Time `json:"r"` //
	PubComp        bool      `json:"c"` // PUBCOMP from remote received

	// tis comes from the packet
	Topic     string `json:"t,omitempty"`
	MessageID uint16 `json:"i,omitempty"` // This is the origin message-ID
	Qos       byte   `json:"q,omitempty"`
	Payload   []byte `json:"d,omitempty"`
}

var packetStoreBucketName string = "packets"
var packetStoreMutex sync.Mutex

// PacketStore will store an packetInfo to persistance store
//
// - need Open() before
//
// - if lastID exist in the store, it will be set to the next available number
//
// - This function is thread-save over mutex
func PacketStore(info PacketInfo, lastID *uint16) (err error) {

	packetStoreMutex.Lock()
	defer packetStoreMutex.Unlock()

	// save it to the db
	err = db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(packetStoreBucketName))
		if b == nil {
			return nil
		}

		// create json-byte-array
		if payload, err := json.Marshal(info); err == nil {

			// try to find a free ID
			newIDBytes := make([]byte, 2)
			for {
				// convert to le-uint16
				binary.LittleEndian.PutUint16(newIDBytes, *lastID)

				existingPacket := b.Get(newIDBytes)
				if existingPacket == nil {
					break
				}

				*lastID++
			}

			// save it to DB
			err = b.Put(newIDBytes, payload)

		}

		return err
	})

	return
}

// PacketIterate will iterate packets
func PacketIterate(iterate func(info PacketInfo, publishPacket *packets.PublishPacket)) (err error) {

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(packetStoreBucketName))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// payload
			var packetInfo PacketInfo

			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				// create a packet
				pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
				pubPacket.MessageID = packetInfo.MessageID
				pubPacket.Retain = packetInfo.OriginClientID == ""
				pubPacket.Dup = packetInfo.OriginClientID != ""
				pubPacket.TopicName = packetInfo.Topic
				pubPacket.Payload = packetInfo.Payload
				pubPacket.Qos = packetInfo.Qos

				// call iterate-function
				iterate(packetInfo, pubPacket)
			}
		}

		return nil
	})

	return
}

func packetGet(clientID *string, topic *string, brokerMessageID *uint16) (found bool, key []byte, packetInfo PacketInfo, err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(packetStoreBucketName))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				found = true

				if clientID != nil {
					found = packetInfo.OriginClientID == *clientID && found
				}

				if topic != nil {
					found = packetInfo.Topic == *topic && found
				}

				if brokerMessageID != nil {
					found = binary.LittleEndian.Uint16(k) == *brokerMessageID && found
				}

				if found == true {
					key = k
					log.Debug().
						Str("bucket", packetStoreBucketName).
						Uint16("key", binary.LittleEndian.Uint16(k)).
						Uint16("mid", packetInfo.MessageID).
						Str("topic", packetInfo.Topic).
						Msg("Delete payload from retained-store")
					break
				}

			}

		}

		return err
	})

	return
}

// PacketExist return if packet exist
//
// If more than one parameter is not NIL all must match
func PacketExist(clientID *string, topic *string, brokerMessageID *uint16) (found bool, err error) {
	found, _, _, err = packetGet(clientID, topic, brokerMessageID)
	return
}

// PacketDelete will delete the first packet it found
//
// - If more than one parameter is not NIL all must match
//
// - This function is thread-save over mutex
func PacketDelete(clientID *string, topic *string, brokerMessageID *uint16) (err error) {

	packetStoreMutex.Lock()
	defer packetStoreMutex.Unlock()

	var found bool
	var foundKey []byte
	var foundPacketInfo PacketInfo

	found, foundKey, foundPacketInfo, err = packetGet(clientID, topic, brokerMessageID)

	if found && err == nil {
		err = db.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte(packetStoreBucketName))
			if b == nil {
				return nil
			}

			if err := b.Delete(foundKey); err != nil {
				log.Debug().
					Str("bucket", packetStoreBucketName).
					Uint16("key", binary.LittleEndian.Uint16(foundKey)).
					Uint16("mid", foundPacketInfo.MessageID).
					Str("topic", foundPacketInfo.Topic).
					Msg("Delete payload from retained-store")
			}
			return err
		})
	}

	return
}

// PacketPubCompSet set the pubcomp-state of the broker-message ID
//
// - this is normal the action when we received a PUBCOMP
//
// - This function is thread-save over mutex
func PacketPubCompSet(brokerMessageID uint16, PubComp bool) (err error) {

	packetStoreMutex.Lock()
	defer packetStoreMutex.Unlock()

	err = db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(packetStoreBucketName))
		if b == nil {
			return nil
		}

		// convert to le-uint16
		brokerMessageIDBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(brokerMessageIDBytes, brokerMessageID)

		if v := b.Get(brokerMessageIDBytes); v != nil {

			var packetInfo PacketInfo
			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				// set PubComp
				packetInfo.PubComp = PubComp

				// create json-byte-array
				var payload []byte
				if payload, err = json.Marshal(packetInfo); err == nil {
					err = b.Put(brokerMessageIDBytes, payload)
				}
			}

		}

		return err
	})

	return
}

// PacketPubCompIsSet return if the original-MessageID receives an pubcomp
//
// This will be normally used, when the origin sender request PUBREL
func PacketPubCompIsSet(MessageID uint16) (pubcomp bool, origMessageID uint16, err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(packetStoreBucketName))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			// payload
			var packetInfo PacketInfo

			// parse it
			if err = json.Unmarshal(v, &packetInfo); err == nil {

				if packetInfo.MessageID == MessageID {
					pubcomp = true
					origMessageID = packetInfo.MessageID
					break
				}

			}

		}

		return err
	})

	return
}
