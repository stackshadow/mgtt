package auth

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/broker"
	"gitlab.com/mgtt/internal/mgtt/clientlist"
)

func onSelfUsernameGet(originClientID string) {
	var err error

	err = clientlist.PublishToClient(
		originClientID,
		"$SYS/self/username",
		[]byte(broker.Current.UserNameOfClient(originClientID)),
	)

	if err != nil {
		log.Error().Err(err).Send()
	}
}
