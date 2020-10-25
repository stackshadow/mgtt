package messagestore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// IterateResendPackets will iterate packages that are stored with StoreResendPacket()
func (store *Store) IterateResendPackets(bucket string, iterate func(storedInfo *StoreResendPacketOption)) (err error) {

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

			//
			storedInfo := StoreResendPacketOption{
				BrokerMessageID: binary.LittleEndian.Uint16(k),
				ClientID:        info.ClientID,
				ResendAt:        info.ResendAt,
				Packet:          publishPacket,
			}

			// call iterate-function
			iterate(&storedInfo)
		}

		return nil
	})

	return
}
