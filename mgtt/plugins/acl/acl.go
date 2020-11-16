package acl

import (
	"strings"

	"github.com/rs/zerolog/log"
	"gitlab.com/mgtt/client"
	"gitlab.com/mgtt/plugin"
)

type aclEntry struct {
	Route string `yaml:"route"`

	// "r" / "w"
	Direction string `yaml:"direction"`

	// true if allow, false if not
	Allow bool `yaml:"allow"`
}

// LocalInit will init the auth-plugin and register it
func LocalInit(ConfigPath string) {

	// OnInit open the config file and watch for changes
	OnInit(ConfigPath)

	newPlugin := plugin.V1{
		OnPublishRequest:          OnPublishRequest,
		OnSendToSubscriberRequest: OnSendToSubscriberRequest,
	}
	plugin.Register("acl", &newPlugin)
}

// OnInit open the config file and watch for changes
func OnInit(ConfigPath string) {
	loadConfig(ConfigPath + "acl.yml")
	go watchConfig()
}

// OnPublishRecvRequest write to broker
func checkACL(clientID string, username string, topic string, direction string) (allowed bool) {
	// if clientID is resend, this is an resended package... we allow this by default
	if clientID == "resend" {
		log.Debug().Str("topic", topic).Msg("This is an resendet packet, we allow it")
		return true
	}

	defer func() {
		if allowed == false {
			log.Warn().Str("topic", topic).Msg("Not allowed")
		}
	}()

	// if username is empty,
	if username == "" {
		username = "_anonym"
	}

	// try to get the acl
	entryArray, exist := config.Rules[username]
	if exist == false {
		return false
	}

	// iterate
	topicArray := strings.Split(topic, "/")
	for _, entry := range entryArray {

		if entry.Direction == direction {

			routeArray := strings.Split(entry.Route, "/")
			if client.Match(routeArray, topicArray) == true {
				allowed = entry.Allow
				return
			}
		}

	}

	return
}

// OnPublishRequest get called when an publisher try to publish to the broker
func OnPublishRequest(clientID string, username string, topic string) (accepted bool) {
	return checkACL(clientID, username, topic, "w")
}

// OnSendToSubscriberRequest get called when the broker try to publish a message to an subscriber
//
// if an plugin set it to false, the message will NOT be sended to clientID
//
// This function gets called BEFORE check if the subscriber subscribe to the topic
//
// clientID: The clientID of the target client
// username: The username of the target client
// publishTopic: The topic the broker try to publish to the subscriber
func OnSendToSubscriberRequest(clientID string, username string, publishTopic string) (accepted bool) {
	return checkACL(clientID, username, publishTopic, "r")
}
