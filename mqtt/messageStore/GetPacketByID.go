package messagestore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

// GetPacketByID return the package by id
func (store *Store) GetPacketByID(bucket string, brokerMessageID uint16) (storedInfo *StoreResendPacketOption, err error) {

	err = store.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		brokerMessageIDBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(brokerMessageIDBytes, brokerMessageID)

		v := b.Get(brokerMessageIDBytes)

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
		storedInfo = &StoreResendPacketOption{
			BrokerMessageID: brokerMessageID,
			ClientID:        info.ClientID,
			ResendAt:        info.ResendAt,
			Packet:          publishPacket,
		}

		return nil
	})

	return
}
