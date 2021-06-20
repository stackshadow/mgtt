package persistance

import (
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
)

// Open will open the DB
func Open(dbFilePath string) (err error) {
	log.Debug().Str("path", dbFilePath).Msg("Try to open DB")
	db, err = bolt.Open(dbFilePath, 0600, nil)

	if err != nil {
		log.Info().Str("path", dbFilePath).Msg("DB opened")
	}

	return
}
