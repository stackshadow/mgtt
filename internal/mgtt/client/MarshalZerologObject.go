package client

import "github.com/rs/zerolog"

func (client *MgttClient) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("clientID", client.id)
}
