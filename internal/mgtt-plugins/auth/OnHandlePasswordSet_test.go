package auth

import (
	"os"
	"testing"
	"time"
)

func TestOnHandlePasswordSet(t *testing.T) {
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
}
