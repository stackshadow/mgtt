package auth

import (
	"os"
	"testing"
	"time"
)

func TestOnHandleMessage(t *testing.T) {

	os.Remove("./integrationtest_auth.yml")

	LocalInit("integrationtest_")
	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/password/set", []byte("admin"))
	time.Sleep(time.Millisecond * 500)

	// check correct password
	if OnAcceptNewClient("integrationtest", "admin", "admin") != true {
		t.Fail()
	}

	// check wrong password
	if OnAcceptNewClient("integrationtest", "admin", "admin2") == true {
		t.Fail()
	}

	// delete user
	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/delete", []byte(""))
	time.Sleep(time.Millisecond * 500)

	// check correct password - but user is missing
	if OnAcceptNewClient("integrationtest", "admin", "admin") == true {
		t.Fail()
	}

	os.Remove("./integrationtest_auth.yml")

}
