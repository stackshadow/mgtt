package acl

import "testing"

func TestPublish(t *testing.T) {

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
}
