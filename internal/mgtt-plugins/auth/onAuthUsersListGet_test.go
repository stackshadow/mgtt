package auth

import (
	"os"
	"testing"
	"time"
)

func TestOnAuthUsersListGet(t *testing.T) {

	os.Remove("./integrationtest_auth.yml")
	LocalInit("integrationtest_")

	// delete user
	OnHandleMessage("integrationtest", "$SYS/auth/users/list", []byte(""))
	time.Sleep(time.Millisecond * 500)

	os.Remove("./integrationtest_auth.yml")

}
