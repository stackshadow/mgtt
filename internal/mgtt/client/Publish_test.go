package client

import (
	"strings"
	"testing"
)

var acceptedRoutes []string = []string{
	"#",
	"/some/#",
	"/some/+/test",
	"/some/+/test",
}

var failedRoutes []string = []string{
	"/some/small/test/that/should",
	"/small/#",
	"/+/dummy/test",
	"/some/+/match",
}

var acceptedTopics []string = []string{
	"/some/small/test/that/should/match",
	"/some/small/test/that/should/match",
	"/some/first/test",
	"/some/second/test",
}

func TestMatch(t *testing.T) {

	for index := range acceptedRoutes {
		acceptedRouteArray := strings.Split(acceptedRoutes[index], "/")
		acceptedTopicsArray := strings.Split(acceptedTopics[index], "/")
		if Match(acceptedRouteArray, acceptedTopicsArray) == false {
			t.Errorf("Route not match topics. %+v == %+v", acceptedRoutes, acceptedTopics)
			t.Fail()
		}
	}

	for index := range failedRoutes {
		acceptedRouteArray := strings.Split(failedRoutes[index], "/")
		acceptedTopicsArray := strings.Split(acceptedTopics[index], "/")
		if Match(acceptedRouteArray, acceptedTopicsArray) == true {
			t.Errorf("Route match topics. %+v == %+v", acceptedRoutes, acceptedTopics)
			t.Fail()
		}
	}

}
