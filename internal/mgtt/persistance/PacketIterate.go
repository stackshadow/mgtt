package persistance

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// PacketIterate will iterate packets
func PacketIterate(bucket string, iterate func(info PacketInfo, publishPacket *packets.PublishPacket)) (err error) {

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

				// create a packet
				pubPacket := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
				pubPacket.MessageID = packetInfo.MessageID
				pubPacket.Retain = false
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
