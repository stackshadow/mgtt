package persistance

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// ClientSessionInfo contains the session infos for an client
type ClientSessionInfo struct {
	ExpireAt      time.Time `json:"e"` //
	Subscriptions []string  `json:"s"` //
}

var clientSessionBucketName string = "sessions"

// SubscriptionsGet will set the subscription of a client-session
func SubscriptionsGet(clientID string) (subscriptions []string) {

	var err error
	var sessionInfo ClientSessionInfo

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(clientSessionBucketName))
		if b == nil {
			return nil
		}

		if v := b.Get([]byte(clientID)); len(v) > 0 {
			// parse json
			err = json.Unmarshal(v, &sessionInfo)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Send()
	}

	return sessionInfo.Subscriptions
}

// SubscriptionsSet will set the subscription of a client-session
func SubscriptionsSet(clientID string, subscriptions []string) {

	var err error
	var sessionInfo ClientSessionInfo

	sessionInfo.Subscriptions = subscriptions

	// save it to the db
	err = db.Update(func(tx *bolt.Tx) error {

		// get bucket
		var b *bolt.Bucket
		b, err = tx.CreateBucketIfNotExists([]byte(clientSessionBucketName))
		if b == nil {
			return nil
		}

		// create json-byte-array
		var payload []byte
		payload, err = json.Marshal(sessionInfo)

		// save it to DB
		err = b.Put([]byte(clientID), payload)

		return err
	})

	if err != nil {
		log.Error().Err(err).Send()
	}
}

// CleanSession will remove the session of an clientID
func CleanSession(clientID string) {
	SubscriptionsSet(clientID, []string{})
}
