package cli

import (
	"os"
)

// AdminTopicsFlag enable all admin-topics
type AdminTopicsFlag bool

// AfterApply setup the json logging
func (v *AdminTopicsFlag) AfterApply() error {

	// ensure, that the env-var is set
	if bool(*v) == true {
		os.Setenv("ENABLE_ADMIN_TOPICS", "true")
	} else {
		os.Setenv("ENABLE_ADMIN_TOPICS", "false")
	}
	return nil
}
