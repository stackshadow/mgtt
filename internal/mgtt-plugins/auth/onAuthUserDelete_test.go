package auth

import (
	"os"
	"testing"
	"time"
)

func TestOnAuthUserDelete(t *testing.T) {

	os.Remove("./integrationtest_auth.yml")
	LocalInit("integrationtest_")

	// delete user
	OnHandleMessage("integrationtest", "$SYS/auth/user/admin/delete", []byte(""))
	time.Sleep(time.Millisecond * 500)

	// check correct password - but user is missing
	if OnAcceptNewClient("integrationtest", "admin", "admin") == true {
		t.Fail()
	}

	os.Remove("./integrationtest_auth.yml")

}
