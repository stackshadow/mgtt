package persistance

import (
	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"gitlab.com/stackshadow/qommunicator/v2/pkg/utils"
)

// Open will open the DB
func MustOpen(dbFilePath string) {
	var err error

	log.Debug().Str("path", dbFilePath).Msg("Try to open DB")
	db, err = bolt.Open(dbFilePath, 0600, nil)
	utils.PanicOnErr(err)

	log.Info().Str("path", dbFilePath).Msg("DB opened")

	return
}
