package auth

import (
	"os"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/internal/mgtt/plugin"
)

// LocalInit will init the auth-plugin and register it
func LocalInit(ConfigPath string) {

	// OnInit open the config file and watch for changes
	OnInit(ConfigPath)

	newPlugin := plugin.V1{
		OnAcceptNewClient: OnAcceptNewClient,
		OnHandleMessage:   OnHandleMessage,
	}
	plugin.Register("auth", &newPlugin)
}

// OnInit open the config file and watch for changes
func OnInit(ConfigPath string) {
	configLoad(ConfigPath + "auth.yml")

	// get username/password from environment
	newUserName := os.Getenv("AUTH_USERNAME")
	newUserPass := os.Getenv("AUTH_PASSWORD")
	os.Unsetenv("AUTH_PASSWORD")

	// get if admin-topics should be enabled
	enableAdminTopics, _ := os.LookupEnv("ENABLE_ADMIN_TOPICS")
	if enableAdminTopics == "true" || enableAdminTopics == "TRUE" || enableAdminTopics == "1" {
		adminTopicsEnabled = true
		log.Info().Msg("admin-topics enabled")
	} else {
		log.Info().Msg("admin-topics disabled")
	}

	// anonymouse is set via environment ?
	_, anonymouse := os.LookupEnv("AUTH_ANONYMOUSE")
	if anonymouse == true {
		config.Anonym = true
	}

	var err error

	if newUserName != "" {
		_, err = userSet(newUserName, &newUserPass, nil)
		if err == nil {
			log.Debug().Str("username", newUserName).Msg("Added new Username")
			err = configSave(filename)
		}
		if err == nil {
			log.Debug().Str("filename", filename).Msg("Config saved")
		}
	}

	if err != nil {
		log.Error().Err(err).Send()
	}

	go configWatch()
}

// OnAcceptNewClient gets called, when a CONNECT-Packet arrived but is not yet added to the list of known clients
func OnAcceptNewClient(clientID string, username string, password string) (accepted bool) {
	return passwordCheck(username, password)
}
