package messagestore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// FindPacket looks for a packet with an clientID and the original messageID
func (store *Store) FindPacket(bucket string, clientID string, originalMessageID uint16) (storedInfo *StoreResendPacketOption, err error) {

	err = store.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			// parse json
			var info packetInfo
			err = json.Unmarshal(v, &info)
			if err != nil {
				return err
			}

			// load the packet
			publishPacketData := bytes.NewBuffer(info.PacketData)
			publishPacketGeneric, err := packets.ReadPacket(publishPacketData)
			if err != nil {
				return err
			}
			publishPacket := publishPacketGeneric.(*packets.PublishPacket)

			if publishPacket.MessageID == originalMessageID && info.ClientID == clientID {
				//
				storedInfo = &StoreResendPacketOption{
					BrokerMessageID: binary.LittleEndian.Uint16(k),
					ClientID:        info.ClientID,
					ResendAt:        info.ResendAt,
					Packet:          publishPacket,
				}

				return nil
			}

		}

		err = errors.New("Not found")

		return err
	})

	return
}
