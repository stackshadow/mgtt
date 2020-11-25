package auth

import (
	"time"

	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"
)

// watchConfig will watch for changes of the config file and reload it when it changes
func watchConfig() (err error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	w.Add(filename)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Info().Str("filename", event.Path).Str("event", event.String()).Msg("File change detected")
				loadConfig(filename)

			case err := <-w.Error:
				log.Error().Err(err).Send()
			case <-w.Closed:
				return
			}
		}
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Error().Err(err).Send()
	}

	return
}
