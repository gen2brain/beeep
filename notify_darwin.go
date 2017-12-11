// +build darwin,!linux,!windows,!js

package beeep

import (
	"os/exec"
)

// Notify sends desktop notification.
func Notify(title, message string) error {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}

	cmd := exec.Command(osa, "-e", `display notification "`+message+`" with title "`+title+`"`)
	return cmd.Start()
}
