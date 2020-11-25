package acl

import (
	"os"
	"testing"
)

func TestPublish(t *testing.T) {

	os.Remove("./integrationtest_auth.yml")
	LocalInit("integrationtest_")

	// we add some acls
	config.Rules["testuser"] = []aclEntry{
		// we not allow write to clients
		aclEntry{
			Direction: "w",
			Route:     "$SYS/broker/clients",
			Allow:     false,
		},

		// and not to the rest
		aclEntry{
			Direction: "w",
			Route:     "$SYS/#",
			Allow:     false,
		},

		// but we allow write to sensors
		aclEntry{
			Direction: "w",
			Route:     "sensors/#",
			Allow:     true,
		},

		// and not to the rest
		aclEntry{
			Direction: "w",
			Route:     "#",
			Allow:     false,
		},
	}

	if OnPublishRequest("", "testuser", "$SYS/broker/clients") == true {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "$SYS/connections/count") == true {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "sensors/temp/first") == false {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "sensors/temp/second") == false {
		t.FailNow()
	}
	if OnPublishRequest("", "testuser", "commands/power/off") == true {
		t.FailNow()
	}

	// we add some acls
	config.Rules["_anonym"] = []aclEntry{
		// we not allow write to clients
		aclEntry{
			Direction: "r",
			Route:     "system/users",
			Allow:     true,
		},
		aclEntry{
			Direction: "w",
			Route:     "system/config",
			Allow:     true,
		},
	}

	if OnPublishRequest("", "", "system/users") == true {
		t.FailNow()
	}
	if OnSendToSubscriberRequest("", "", "system/users") == false {
		t.FailNow()
	}

	if OnPublishRequest("", "", "system/config") == false {
		t.FailNow()
	}
	if OnSendToSubscriberRequest("", "", "system/config") == true {
		t.FailNow()
	}

	os.Remove("./integrationtest_auth.yml")
}
