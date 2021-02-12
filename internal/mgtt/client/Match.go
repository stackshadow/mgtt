package client

import "strings"

// MatchRoute compare an route with an topic
func MatchRoute(route string, topic string) (match bool) {

	routeArray := strings.Split(route, "/")
	topicArray := strings.Split(topic, "/")

	return Match(routeArray, topicArray)
}

// Match compare an selector route with an given topic
func Match(route []string, topic []string) bool {
	if len(route) == 0 {
		return len(topic) == 0
	}

	if len(topic) == 0 {
		return route[0] == "#"
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return Match(route[1:], topic[1:])
	}
	return false
}
